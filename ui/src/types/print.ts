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

export const QUICK_DENY_REASONS = [
    "File contains inappropriate content",
    "Model too large for available printers",
    "Poor file quality/corrupted",
    "Duplicate request",
    "Against printing policy",
    "Insufficient detail for printing",
    "Copyright infringement",
    "Safety concerns"
];