import { Request } from "express";
import { ShareInfoResponse, FolderListResponse, SaveFileParams, SaveFileResult, RenameFileResult } from "./cloud"; // Added SaveFileResult and RenameFileResult

export interface ICloudStorageService {
  setCookie(req: Request): Promise<void>;
  getShareInfo(shareCode: string, receiveCode?: string): Promise<ShareInfoResponse>;
  getFolderList(parentCid?: string): Promise<FolderListResponse>;
  saveSharedFile(params: SaveFileParams): Promise<SaveFileResult>; // Changed return type
  renameFile(fileId: string, targetFileName: string, folderId?: string): Promise<RenameFileResult>;
}
