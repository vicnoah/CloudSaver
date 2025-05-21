import { AxiosInstance, AxiosHeaders } from "axios";
import { logger } from "../utils/logger";
import { createAxiosInstance } from "../utils/axiosInstance";
import { injectable } from "inversify";
import { Request } from "express";
import UserSetting from "../models/UserSetting";
import {
  ShareInfoResponse,
  FolderListResponse,
  QuarkFolderItem,
  SaveFileParams,
  SaveFileResult,
  RenameFileResult,
} from "../types/cloud";
import { ICloudStorageService } from "@/types/services";

interface QuarkShareInfo {
  stoken?: string;
  pwdId?: string;
  fileSize?: number;
  list: {
    fid: string;
    file_name: string;
    file_type: number;
    share_fid_token: string;
  }[];
}

@injectable()
export class QuarkService implements ICloudStorageService {
  private api: AxiosInstance;
  private cookie: string = "";

  constructor() {
    this.api = createAxiosInstance(
      "https://drive-h.quark.cn",
      AxiosHeaders.from({
        cookie: this.cookie,
        accept: "application/json, text/plain, */*",
        "accept-language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
        "content-type": "application/json",
        priority: "u=1, i",
        "sec-ch-ua": '"Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="24"',
        "sec-ch-ua-mobile": "?0",
        "sec-ch-ua-platform": '"Windows"',
        "sec-fetch-dest": "empty",
        "sec-fetch-mode": "cors",
        "sec-fetch-site": "same-site",
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
    if (userSetting && userSetting.dataValues.quarkCookie) {
      this.cookie = userSetting.dataValues.quarkCookie;
    } else {
      throw new Error("请先设置夸克网盘cookie");
    }
  }

  async getShareInfo(pwdId: string, passcode = ""): Promise<ShareInfoResponse> {
    const response = await this.api.post(
      `/1/clouddrive/share/sharepage/token?pr=ucpro&fr=pc&uc_param_str=&__dt=994&__t=${Date.now()}`,
      {
        pwd_id: pwdId,
        passcode,
      }
    );
    if (response.data?.status === 200 && response.data.data) {
      const fileInfo = response.data.data;
      if (fileInfo.stoken) {
        const res = await this.getShareList(pwdId, fileInfo.stoken);
        return {
          data: res,
        };
      }
    }
    throw new Error("获取夸克分享信息失败");
  }

  async getShareList(pwdId: string, stoken: string): Promise<ShareInfoResponse["data"]> {
    const response = await this.api.get("/1/clouddrive/share/sharepage/detail", {
      params: {
        pr: "ucpro",
        fr: "pc",
        uc_param_str: "",
        pwd_id: pwdId,
        stoken: stoken,
        pdir_fid: "0",
        force: "0",
        _page: "1",
        _size: "50",
        _fetch_banner: "1",
        _fetch_share: "1",
        _fetch_total: "1",
        _sort: "file_type:asc,updated_at:desc",
        __dt: "1589",
        __t: Date.now(),
      },
    });
    if (response.data?.data) {
      const list = response.data.data.list
        .filter((item: QuarkShareInfo["list"][0]) => item.fid)
        .map((folder: QuarkShareInfo["list"][0]) => ({
          fileId: folder.fid,
          fileName: folder.file_name,
          fileIdToken: folder.share_fid_token,
        }));
      return {
        list,
        pwdId,
        stoken,
        fileSize: response.data.data.share?.size || 0,
      };
    } else {
      return {
        list: [],
      };
    }
  }

  async getFolderList(parentCid = "0"): Promise<FolderListResponse> {
    const response = await this.api.get("/1/clouddrive/file/sort", {
      params: {
        pr: "ucpro",
        fr: "pc",
        uc_param_str: "",
        pdir_fid: parentCid,
        _page: "1",
        _size: "100",
        _fetch_total: "false",
        _fetch_sub_dirs: "1",
        _sort: "",
        __dt: "2093126",
        __t: Date.now(),
      },
    });
    if (response.data?.data && response.data.data.list) {
      const data = response.data.data.list
        .filter((item: QuarkFolderItem) => item.fid && item.file_type === 0)
        .map((folder: QuarkFolderItem) => ({
          cid: folder.fid,
          name: folder.file_name,
          path: [],
        }));
      return {
        data,
      };
    } else {
      const message = "获取夸克目录列表失败:" + response.data.error;
      logger.error(message);
      throw new Error(message);
    }
  }

  async saveSharedFile(params: SaveFileParams): Promise<SaveFileResult> {
    const quarkParams = {
      fid_list: params.fids,
      fid_token_list: params.fidTokens,
      to_pdir_fid: params.folderId,
      pwd_id: params.shareCode,
      stoken: params.receiveCode, // Assuming receiveCode for Quark is stoken
      pdir_fid: "0", // Default or make configurable if needed
      scene: "link", // Default or make configurable if needed
    };

    try {
      const response = await this.api.post(
        `/1/clouddrive/share/sharepage/save?pr=ucpro&fr=pc&uc_param_str=&__dt=208097&__t=${Date.now()}`,
        quarkParams
      );

      // Assuming a successful HTTP call means the operation was accepted by Quark.
      // Quark API success is usually indicated by HTTP status 200 and specific fields in response.data
      // response.data.message seems to be the primary indicator from the original code.
      let finalMessage = response.data.message || "文件已成功提交保存";
      let saveSuccess = true; // Assuming the api call itself implies save initiation success

      // SPECULATIVE: Extract new file ID and name from Quark's response.
      // Actual fields in response.data.data for the new file's ID and name are unknown.
      // These are placeholders. `params.fids[0]` is the ID of the file *being shared*.
      // We need the ID of the file *after* it's saved to the user's drive.
      const savedFileId = response.data.data?.fid || response.data.data?.new_fid || response.data.data?.file_info_list?.[0]?.fid; // Example placeholder
      const currentName = response.data.data?.file_name || response.data.data?.new_name || response.data.data?.file_info_list?.[0]?.file_name; // Example placeholder

      let resultFileId = savedFileId; 
      let resultActualFileName = currentName;

      if (params.targetFileName && savedFileId) {
        if (currentName && params.targetFileName.trim() === currentName.trim()) {
          logger.info(`Quark: Target name "${params.targetFileName}" is same as current name "${currentName}". Skipping rename.`);
          finalMessage += `. Filename is already as requested: ${currentName}`;
        } else if (!currentName) {
          logger.warn(`Quark: Current name for file ${savedFileId} could not be determined. Proceeding with rename attempt to ${params.targetFileName}.`);
          const renameResult = await this.renameFile(savedFileId, params.targetFileName, params.folderId);
          if (renameResult.success) {
            finalMessage += `. Successfully renamed to ${params.targetFileName}`;
            resultActualFileName = params.targetFileName;
          } else {
            finalMessage += `. Failed to rename: ${renameResult.message}`;
          }
        } else { // currentName is known and different
          logger.info(`Quark: Attempting to rename file ${savedFileId} (current: "${currentName}") to ${params.targetFileName}`);
          const renameResult = await this.renameFile(savedFileId, params.targetFileName, params.folderId);
          if (renameResult.success) {
            finalMessage += `. Successfully renamed to ${params.targetFileName}`;
            resultActualFileName = params.targetFileName;
          } else {
            finalMessage += `. Failed to rename: ${renameResult.message}`;
          }
        }
      } else if (params.targetFileName) {
        logger.warn("Quark: targetFileName provided, but could not determine savedFileId from response. Skipping rename.");
        finalMessage += ". Could not rename: missing file ID from save response.";
      }

      return {
        success: saveSuccess,
        message: finalMessage,
        data: response.data.data,
        fileId: resultFileId,
        actualFileName: resultActualFileName,
      };
    } catch (error: any) { // Catch block to handle axios errors or other exceptions
      logger.error("保存Quark文件请求异常:", error);
      const errorMessage = error.response?.data?.message || (error instanceof Error ? error.message : "未知错误");
      return {
        success: false,
        message: "保存Quark文件请求失败: " + errorMessage,
        data: error.response?.data?.data,
      };
    }
  }

  async renameFile(fileId: string, targetFileName: string, folderId?: string): Promise<RenameFileResult> {
    logger.warn(`renameFile called for Quark service (fileId: ${fileId}, targetName: ${targetFileName}, folderId: ${folderId}), but it is not yet implemented.`);
    return {
      success: false,
      message: "File renaming for Quark service is not yet implemented.",
    };
  }
}
