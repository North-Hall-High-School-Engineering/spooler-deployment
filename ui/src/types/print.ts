export type Print = {
    ID: number;
    Status: string;
    UploadedFileName: string;
    StoredFileName: string;
    RequestedFilamentColor: string;
    CreatedAt: string;
    DenialReason?: string;
};

export type PrintStatus = "approval_pending" | "pending_print" | "printing" | "denied" | "completed" | "failed" | "canceled" | "paused";