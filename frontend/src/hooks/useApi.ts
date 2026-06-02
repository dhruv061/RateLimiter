import { useEffect, useState } from "react";
import { api } from "../services/api";
import { useGlobalFilter } from "../context/GlobalFilterContext";

export function useApi<T>(path: string, fallback: T) {
  const [data, setData] = useState<T>(fallback);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Access global filter state to reactively reload when selections change
  const { selectedDomain, selectedRange, customRange, refreshTrigger } = useGlobalFilter();

  useEffect(() => {
    let active = true;
    setLoading(true);
    api<T>(path)
      .then((result) => {
        if (active) {
          setData(result);
          setError(null);
        }
      })
      .catch((err: Error) => {
        if (active) {
          setError(err.message);
        }
      })
      .finally(() => {
        if (active) {
          setLoading(false);
        }
      });
    return () => {
      active = false;
    };
  }, [path, selectedDomain, selectedRange, customRange.start, customRange.end, refreshTrigger]);

  return { data, setData, loading, error };
}
