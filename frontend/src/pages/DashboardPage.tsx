import { useState, useEffect } from "react";
import { useAuthStore } from "../stores/authStore";

interface AnalyticsSummary {
  totalEntries: number;
  bristolDistribution: { type: number; count: number }[];
  recentEntries: {
    id: string;
    bristolType: number;
    createdAt: string;
    satisfaction?: number;
  }[];
  averageSatisfaction?: number;
}

interface EntryResponse {
  id: string;
  bristolType: number;
  volume?: string;
  color?: string;
  notes?: string;
  createdAt: string;
  userId: string;
}

interface EntriesApiResponse {
  entries: EntryResponse[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    pages: number;
  };
}

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3002";

const getBristolTypeDescription = (type: number): string => {
  const descriptions = {
    1: "Hard lumps",
    2: "Lumpy sausage",
    3: "Cracked sausage",
    4: "Smooth sausage",
    5: "Soft blobs",
    6: "Fluffy pieces",
    7: "Watery"
  };
  return descriptions[type as keyof typeof descriptions] || "Unknown";
};

const getBristolTypeCategory = (type: number): string => {
  if (type <= 2) return "Constipated";
  if (type <= 4) return "Normal";
  return "Loose";
};

const getThisWeekCount = (entries: EntryResponse[]): number => {
  const oneWeekAgo = new Date();
  oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);

  return entries.filter((entry) => new Date(entry.createdAt) >= oneWeekAgo)
    .length;
};

const getAverageBristolType = (
  bristolDistribution: { type: number; count: number }[]
): number => {
  if (bristolDistribution.length === 0) return 0;

  const totalCount = bristolDistribution.reduce(
    (sum, item) => sum + item.count,
    0
  );
  const weightedSum = bristolDistribution.reduce(
    (sum, item) => sum + item.type * item.count,
    0
  );

  return totalCount > 0 ? Number((weightedSum / totalCount).toFixed(1)) : 0;
};

export function DashboardPage() {
  const { token } = useAuthStore();
  const [analytics, setAnalytics] = useState<AnalyticsSummary | null>(null);
  const [recentEntries, setRecentEntries] = useState<EntryResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string>("");

  const fetchAnalytics = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/analytics/summary`, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error("Failed to fetch analytics");
      }

      const data: AnalyticsSummary = await response.json();
      setAnalytics(data);
    } catch (error) {
      console.error("Error fetching analytics:", error);
      setError(
        error instanceof Error ? error.message : "Failed to load analytics"
      );
    }
  };

  const fetchRecentEntries = async () => {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/entries?limit=5&sortOrder=desc`,
        {
          headers: {
            Authorization: `Bearer ${token}`
          }
        }
      );

      if (!response.ok) {
        throw new Error("Failed to fetch entries");
      }

      const data: EntriesApiResponse = await response.json();
      setRecentEntries(data.entries);
    } catch (error) {
      console.error("Error fetching recent entries:", error);
      setError(
        error instanceof Error ? error.message : "Failed to load recent entries"
      );
    }
  };

  useEffect(() => {
    const loadDashboardData = async () => {
      setIsLoading(true);
      setError("");

      if (token) {
        await Promise.all([fetchAnalytics(), fetchRecentEntries()]);
      }

      setIsLoading(false);
    };

    loadDashboardData();
  }, [token]);

  if (isLoading) {
    return (
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">💩 Dashboard</h1>
        <div className="card text-center py-8">
          <p className="text-gray-600">Loading your poo data...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">💩 Dashboard</h1>
        <div className="card text-center py-8">
          <p className="text-red-600">Error: {error}</p>
        </div>
      </div>
    );
  }

  const averageBristol = analytics
    ? getAverageBristolType(analytics.bristolDistribution)
    : 0;
  const thisWeekCount = getThisWeekCount(recentEntries);

  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">💩 Dashboard</h1>

      <div className="grid md:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">Total Entries</h3>
          <p className="text-3xl font-bold text-poo-brown-600">
            {analytics?.totalEntries || 0}
          </p>
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold mb-2">Avg Bristol Score</h3>
          <p className="text-3xl font-bold text-poo-brown-600">
            {averageBristol || "N/A"}
          </p>
          {averageBristol > 0 && (
            <p className="text-sm text-gray-600 mt-1">
              {getBristolTypeCategory(averageBristol)}
            </p>
          )}
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold mb-2">This Week</h3>
          <p className="text-3xl font-bold text-poo-brown-600">
            {thisWeekCount}
          </p>
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">Recent Entries</h3>
          {recentEntries.length === 0 ? (
            <p className="text-gray-600">
              No entries yet. Start tracking your bowel movements!
            </p>
          ) : (
            <div className="space-y-3">
              {recentEntries.map((entry) => (
                <div
                  key={entry.id}
                  className="flex justify-between items-center p-3 bg-gray-50 rounded"
                >
                  <div>
                    <div className="font-medium">
                      Bristol Type {entry.bristolType}
                    </div>
                    <div className="text-sm text-gray-600">
                      {getBristolTypeDescription(entry.bristolType)}
                    </div>
                    {entry.volume && (
                      <div className="text-xs text-gray-500">
                        Volume: {entry.volume}
                      </div>
                    )}
                  </div>
                  <div className="text-right">
                    <div className="text-sm text-gray-600">
                      {new Date(entry.createdAt).toLocaleDateString()}
                    </div>
                    <div className="text-xs text-gray-500">
                      {new Date(entry.createdAt).toLocaleTimeString([], {
                        hour: "2-digit",
                        minute: "2-digit"
                      })}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold mb-4">Quick Stats</h3>
          {analytics ? (
            <div className="space-y-4">
              <div>
                <h4 className="font-medium mb-2">Bristol Type Distribution</h4>
                <div className="space-y-2">
                  {analytics.bristolDistribution.map((item) => (
                    <div
                      key={item.type}
                      className="flex justify-between items-center"
                    >
                      <span className="text-sm">
                        Type {item.type}: {getBristolTypeDescription(item.type)}
                      </span>
                      <div className="flex items-center">
                        <span className="text-sm font-medium mr-2">
                          {item.count}
                        </span>
                        <div className="w-16 h-2 bg-gray-200 rounded">
                          <div
                            className="h-2 bg-poo-brown-500 rounded"
                            style={{
                              width: `${(item.count / (analytics.totalEntries || 1)) * 100}%`
                            }}
                          ></div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {analytics.averageSatisfaction && (
                <div>
                  <h4 className="font-medium mb-2">Average Satisfaction</h4>
                  <div className="flex items-center">
                    <span className="text-2xl font-bold text-poo-brown-600 mr-2">
                      {analytics.averageSatisfaction.toFixed(1)}
                    </span>
                    <span className="text-gray-600">/10</span>
                  </div>
                </div>
              )}

              <div>
                <h4 className="font-medium mb-2">Health Trend</h4>
                <div className="text-sm text-gray-600">
                  {averageBristol >= 3 && averageBristol <= 4
                    ? "🟢 Your bowel movements are in the normal range!"
                    : averageBristol < 3
                      ? "🟡 You might be experiencing constipation. Consider more fiber."
                      : "🟡 Your stools are on the loose side. Monitor your diet."}
                </div>
              </div>
            </div>
          ) : (
            <p className="text-gray-600">
              Analytics will appear here once you start tracking...
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
