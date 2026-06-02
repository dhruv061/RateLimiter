import React, { createContext, useContext, useState, useEffect } from "react";
import { api, currentFilters } from "../services/api";

export interface Domain {
  id: number;
  domain_name: string;
  access_log_path: string;
  error_log_path: string;
  blocked_ip_file_path: string;
  fail2ban_jail_name: string;
  server_name: string;
  description: string;
  is_valid: boolean;
  last_validated_at?: string;
  created_at: string;
  updated_at: string;
}

interface CustomRange {
  start: Date | null;
  end: Date | null;
}

interface GlobalFilterContextType {
  selectedDomain: number;
  setSelectedDomain: (id: number) => void;
  selectedRange: string;
  setSelectedRange: (range: string) => void;
  customRange: CustomRange;
  setCustomRange: (range: CustomRange) => void;
  domains: Domain[];
  loadingDomains: boolean;
  refreshDomains: () => Promise<void>;
  refreshTrigger: number;
  triggerRefresh: () => void;
}

const GlobalFilterContext = createContext<GlobalFilterContextType | undefined>(undefined);

export const useGlobalFilter = () => {
  const context = useContext(GlobalFilterContext);
  if (!context) {
    throw new Error("useGlobalFilter must be used within a GlobalFilterProvider");
  }
  return context;
};

export const GlobalFilterProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [selectedDomain, setSelectedDomainState] = useState<number>(() => {
    const saved = localStorage.getItem("global-filter-domain");
    return saved ? Number(saved) : 0;
  });

  const [selectedRange, setSelectedRangeState] = useState<string>(() => {
    return localStorage.getItem("global-filter-range") || "last_24h";
  });

  const [customRange, setCustomRangeState] = useState<CustomRange>({
    start: null,
    end: null,
  });

  const [domains, setDomains] = useState<Domain[]>([]);
  const [loadingDomains, setLoadingDomains] = useState(true);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const triggerRefresh = () => setRefreshTrigger((prev) => prev + 1);

  const setSelectedDomain = (id: number) => {
    setSelectedDomainState(id);
    localStorage.setItem("global-filter-domain", String(id));
  };

  const setSelectedRange = (range: string) => {
    setSelectedRangeState(range);
    localStorage.setItem("global-filter-range", range);
  };

  const setCustomRange = (range: CustomRange) => {
    setCustomRangeState(range);
  };

  const refreshDomains = async () => {
    try {
      setLoadingDomains(true);
      const list = await api<Domain[]>("/api/domains");
      setDomains(list);
    } catch (err) {
      console.error("Failed to load domains:", err);
    } finally {
      setLoadingDomains(false);
    }
  };

  useEffect(() => {
    refreshDomains();
  }, []);

  // Synchronize currentFilters in api.ts whenever React state changes
  useEffect(() => {
    currentFilters.domainId = selectedDomain;

    const now = new Date();
    let startTime = "";
    let endTime = "";

    if (selectedRange === "last_1h") {
      startTime = new Date(now.getTime() - 60 * 60 * 1000).toISOString();
      endTime = now.toISOString();
    } else if (selectedRange === "last_24h") {
      startTime = new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString();
      endTime = now.toISOString();
    } else if (selectedRange === "last_7d") {
      startTime = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000).toISOString();
      endTime = now.toISOString();
    } else if (selectedRange === "last_30d") {
      startTime = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000).toISOString();
      endTime = now.toISOString();
    } else if (selectedRange === "custom" && customRange.start && customRange.end) {
      startTime = customRange.start.toISOString();
      endTime = customRange.end.toISOString();
    }

    currentFilters.startTime = startTime;
    currentFilters.endTime = endTime;
  }, [selectedDomain, selectedRange, customRange]);

  return (
    <GlobalFilterContext.Provider
      value={{
        selectedDomain,
        setSelectedDomain,
        selectedRange,
        setSelectedRange,
        customRange,
        setCustomRange,
        domains,
        loadingDomains,
        refreshDomains,
        refreshTrigger,
        triggerRefresh,
      }}
    >
      {children}
    </GlobalFilterContext.Provider>
  );
};
