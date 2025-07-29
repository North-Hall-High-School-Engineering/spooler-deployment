export type PrintStatus =
  | "approval_pending"
  | "pending_print"
  | "printing"
  | "denied"
  | "completed"
  | "failed"
  | "canceled"
  | "paused";

export interface Print {
  ID: number;
  UserID: number;
  Status: PrintStatus;
  Progress: number;
  UploadedFileName: string;
  StoredFileName: string;
  RequestedFilamentColor: string;
  DenialReason?: string;
  CreatedAt: string;
  UpdatedAt: string;
}