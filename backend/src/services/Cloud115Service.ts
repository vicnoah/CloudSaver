import { AxiosHeaders, AxiosInstance } from "axios"; // 导入 AxiosHeaders
import { createAxiosInstance } from "../utils/axiosInstance";
import { ShareInfoResponse, FolderListResponse, SaveFileParams, SaveFileResult, RenameFileResult } from "../types/cloud";
import { injectable } from "inversify";
import { Request } from "express";
import UserSetting from "../models/UserSetting";
import { ICloudStorageService } from "@/types/services";
import { logger } from "../utils/logger";

interface Cloud115ListItem {
  cid: string;
  n: string;
  s: number;
}

interface Cloud115FolderItem {
  cid: string;
  n: string;
  ns: number;
}

@injectable()
export class Cloud115Service implements ICloudStorageService {
  private api: AxiosInstance;
  private cookie: string = "";

  constructor() {
    this.api = createAxiosInstance(
      "https://webapi.115.com",
      AxiosHeaders.from({
        Host: "webapi.115.com",
        Connection: "keep-alive",
        xweb_xhr: "1",
        Origin: "",
        "Content-Type": "application/x-www-form-urlencoded",
        "User-Agent":
          "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 MicroMessenger/6.8.0(0x16080000) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.9(0x13080910) XWEB/1227",
        Accept: "*/*",
        "Sec-Fetch-Site": "cross-site",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Dest": "empty",
        Referer: "https://servicewechat.com/wx2c744c010a61b0fa/94/page-frame.html",
        "Accept-Encoding": "gzip, deflate, br",
        "Accept-Language": "zh-CN,zh;q=0.9",
      })
    );

    this.api.interceptors.request.use((config) => {
      config.headers.cookie = this.cookie;
      return config;
    });
  }

  async setCookie(req: Request): Promise<void> {
    const userId = req.user?.userId;
    const userSetting = await UserSetting.findOne({
      where: { userId },
    });
    if (userSetting && userSetting.dataValues.cloud115Cookie) {
      this.cookie = userSetting.dataValues.cloud115Cookie;
    } else {
      throw new Error("请先设置115网盘cookie");
    }
  }

  async getShareInfo(shareCode: string, receiveCode = ""): Promise<ShareInfoResponse> {
    const response = await this.api.get("/share/snap", {
      params: {
        share_code: shareCode,
        receive_code: receiveCode,
        offset: 0,
        limit: 20,
        cid: "",
      },
    });
    if (response.data?.state && response.data.data?.list?.length > 0) {
      return {
        data: {
          list: response.data.data.list.map((item: Cloud115ListItem) => ({
            fileId: item.cid,
            fileName: item.n,
            fileSize: item.s,
          })),
        },
      };
    } else {
      logger.error("未找到文件信息:", response.data);
      throw new Error("未找到文件信息");
    }
  }

  async getFolderList(parentCid = "0"): Promise<FolderListResponse> {
    const response = await this.api.get("/files", {
      params: {
        aid: 1,
        cid: parentCid,
        o: "user_ptime",
        asc: 1,
        offset: 0,
        show_dir: 1,
        limit: 50,
        type: 0,
        format: "json",
        star: 0,
        suffix: "",
        natsort: 0,
        snap: 0,
        record_open_time: 1,
        fc_mix: 0,
      },
    });
    if (response.data?.state) {
      return {
        data: response.data.data
          .filter((item: Cloud115FolderItem) => item.cid && !!item.ns)
          .map((folder: Cloud115FolderItem) => ({
            cid: folder.cid,
            name: folder.n,
            path: response.data.path,
          })),
      };
    } else {
      logger.error("获取目录列表失败:", response.data.error);
      throw new Error("获取115pan目录列表失败:" + response.data.error);
    }
  }

  async saveSharedFile(params: SaveFileParams): Promise<SaveFileResult> {
    const param = new URLSearchParams({
      cid: params.folderId || "",
      share_code: params.shareCode || "",
      receive_code: params.receiveCode || "",
      file_id: params.fids?.[0] || "",
    });
    try {
      const response = await this.api.post("/share/receive", param.toString());
      logger.info("保存文件:", response.data);
      if (response.data.state) {
        let finalMessage = response.data.error || "文件保存成功";
        let saveSuccess = true;

        // SPECULATIVE: Extract new file ID and name from response.
        // Actual fields in response.data.data for the new file's ID and name are unknown for 115's /share/receive.
        // These are placeholders. If these details aren't in the response, renaming accurately is not possible
        // without further API calls to find the file.
        const savedFileId = response.data.data?.new_cid || response.data.data?.cid || params.fids?.[0]; // Example placeholder for file ID
        const currentName = response.data.data?.new_name || response.data.data?.name; // Example placeholder for file name

        // Store initial save details, especially if fileId or actualFileName can be confirmed from response.data.data
        let resultFileId = savedFileId; // Or a more confirmed ID from response.data.data
        let resultActualFileName = currentName; // Or a more confirmed name

        if (params.targetFileName && savedFileId) {
          if (currentName && params.targetFileName.trim() === currentName.trim()) {
            logger.info(`Cloud115: Target name "${params.targetFileName}" is same as current name "${currentName}". Skipping rename.`);
            finalMessage += `. Filename is already as requested: ${currentName}`;
          } else if (!currentName) {
            logger.warn(`Cloud115: Current name for file ${savedFileId} could not be determined. Proceeding with rename attempt to ${params.targetFileName}.`);
            const renameResult = await this.renameFile(savedFileId, params.targetFileName, params.folderId);
            if (renameResult.success) {
              finalMessage += `. Successfully renamed to ${params.targetFileName}`;
              resultActualFileName = params.targetFileName;
            } else {
              finalMessage += `. Failed to rename: ${renameResult.message}`;
            }
          } else { // currentName is known and different
            logger.info(`Cloud115: Attempting to rename file ${savedFileId} (current: "${currentName}") to ${params.targetFileName}`);
            const renameResult = await this.renameFile(savedFileId, params.targetFileName, params.folderId);
            if (renameResult.success) {
              finalMessage += `. Successfully renamed to ${params.targetFileName}`;
              resultActualFileName = params.targetFileName;
            } else {
              finalMessage += `. Failed to rename: ${renameResult.message}`;
            }
          }
        } else if (params.targetFileName) {
          logger.warn("Cloud115: targetFileName provided, but could not determine savedFileId. Skipping rename.");
          finalMessage += ". Could not rename: missing file ID from save response.";
        }

        return {
          success: saveSuccess,
          message: finalMessage,
          data: response.data.data,
          fileId: resultFileId, // Populate with the determined file ID
          actualFileName: resultActualFileName, // Populate with the determined actual file name
        };
      } else {
        logger.error("保存文件失败:", response.data.error);
        return {
          success: false,
          message: "保存115pan文件失败:" + response.data.error,
          data: response.data.data,
        };
      }
    } catch (error) {
      logger.error("保存115pan文件请求异常:", error);
      // Ensure even network or unexpected errors return SaveFileResult
      const errorMessage = error instanceof Error ? error.message : "未知错误";
      return {
        success: false,
        message: "保存115pan文件请求异常: " + errorMessage,
      };
    }
  }

  async renameFile(fileId: string, targetFileName: string, folderId?: string): Promise<RenameFileResult> {
    logger.warn(`renameFile called for 115 service (fileId: ${fileId}, targetName: ${targetFileName}, folderId: ${folderId}), but it is not yet implemented.`);
    return {
      success: false,
      message: "File renaming for 115 service is not yet implemented.",
    };
  }
}
