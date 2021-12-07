export interface RequestParams {
  endpoint: string;
  body?: any;
  verb?: string;
  authNeeded?: boolean;
}

export interface PaymentURL {
    URL: string;
}

export interface User {
    id: string;
    name: string;
    count: number;
}