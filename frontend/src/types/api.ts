export type ApiEnvelope<T> = {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
};

export type User = {
  id: number;
  username: string;
  email: string;
  role: string;
  must_change_pass: boolean;
};

export type DashboardStats = {
  total_requests_today: number;
  requests_last_hour: number;
  current_rps: number;
  peak_rps: number;
  total_429_today: number;
  count_429_last_hour: number;
  current_429_rate: number;
  total_bans_today: number;
  active_bans: number;
  bans_24h: number;
  unbans_today: number;
  nginx_status: string;
  fail2ban_status: string;
  database_status: string;
  service_status: string;
};

export type Ban = {
  id: number;
  ip_address: string;
  country: string;
  country_code: string;
  region: string;
  city: string;
  asn: string;
  isp: string;
  jail: string;
  reason: string;
  ban_time: string;
  unban_time?: string;
  ban_duration: number;
  request_count: number;
  violation_count: number;
  is_active: boolean;
  created_at: string;
};

export type PaginatedResponse<T> = {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
};

export type TrafficStat = {
  id: number;
  timestamp: string;
  total_requests: number;
  unique_ips: number;
  status_429: number;
  status_403: number;
  avg_response_time: number;
  period: string;
};

export type TopOffender = {
  ip_address: string;
  country: string;
  country_code: string;
  total_requests: number;
  violation_count: number;
  ban_count: number;
};

export type CountryStats = {
  country: string;
  country_code: string;
  requests: number;
  violations: number;
  bans: number;
};

export type WhitelistEntry = {
  id: number;
  ip_address: string;
  description: string;
  added_by: string;
  created_at: string;
};

export type AuditLog = {
  id: number;
  user_id: number;
  username: string;
  action: string;
  target: string;
  details: string;
  ip_address: string;
  created_at: string;
};

export type LiveRequest = {
  timestamp: string;
  ip_address: string;
  method: string;
  url: string;
  status_code: number;
  response_time: number;
  user_agent: string;
  bytes_sent: number;
};
