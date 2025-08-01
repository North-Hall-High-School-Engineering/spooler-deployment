export type User = {
    ID: number;
    Email: string;
    FirstName: string;
    LastName: string;
    Role: "user" | "admin" | "officer";
    CreatedAt: string;
    UpdatedAt: string;
}