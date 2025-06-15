import { useState, useEffect } from "react";
import { useAuthStore } from "../stores/authStore";
import Logo from "../components/Logo";

interface StoolEntry {
  bristolType: number;
  volume?: string;
  color?: string;
  notes?: string;
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
    1: "Hard lumps (Severe constipation)",
    2: "Lumpy sausage (Mild constipation)",
    3: "Cracked sausage (Normal)",
    4: "Smooth sausage (Ideal)",
    5: "Soft blobs (Lacking fiber)",
    6: "Fluffy pieces (Mild diarrhea)",
    7: "Watery (Severe diarrhea)"
  };
  return descriptions[type as keyof typeof descriptions] || "Unknown";
};

export function NewEntryPage() {
  const { token } = useAuthStore();
  const [formData, setFormData] = useState<StoolEntry>({
    bristolType: 0,
    volume: "",
    color: "",
    notes: ""
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] = useState<
    "idle" | "success" | "error"
  >("idle");
  const [errorMessage, setErrorMessage] = useState<string>("");
  const [entries, setEntries] = useState<EntryResponse[]>([]);
  const [isLoadingEntries, setIsLoadingEntries] = useState(true);

  const fetchEntries = async () => {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/entries?limit=10&sortOrder=desc`,
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
      setEntries(data.entries);
    } catch (error) {
      console.error("Error fetching entries:", error);
    } finally {
      setIsLoadingEntries(false);
    }
  };

  useEffect(() => {
    if (token) {
      fetchEntries();
    }
  }, [token]);

  const handleInputChange = (
    field: keyof StoolEntry,
    value: string | number
  ) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.bristolType === 0) {
      setErrorMessage("Please select a Bristol Type");
      setSubmitStatus("error");
      return;
    }

    setIsSubmitting(true);
    setSubmitStatus("idle");
    setErrorMessage("");

    try {
      const submitData = {
        bristolType: formData.bristolType,
        volume: formData.volume || undefined,
        color: formData.color || undefined,
        notes: formData.notes || undefined
      };

      const response = await fetch(`${API_BASE_URL}/api/entries`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(submitData)
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to save entry");
      }

      const savedEntry: EntryResponse = await response.json();
      console.log("Entry saved successfully:", savedEntry);

      setSubmitStatus("success");
      setFormData({
        bristolType: 0,
        volume: "",
        color: "",
        notes: ""
      });

      // Refresh the entries list to show the new entry
      await fetchEntries();
    } catch (error) {
      console.error("Error saving entry:", error);
      setErrorMessage(
        error instanceof Error ? error.message : "An error occurred"
      );
      setSubmitStatus("error");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Log New Entry</h1>

      <div className="card">
        <p className="text-center text-gray-600 mb-8 flex items-center justify-center gap-2">
          Time to document another masterpiece! <Logo size={24} />
        </p>

        {submitStatus === "success" && (
          <div className="mb-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
            Entry saved successfully! Keep tracking your progress.
          </div>
        )}

        {submitStatus === "error" && (
          <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {errorMessage}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Bristol Stool Chart Type (1-7) *
            </label>
            <select
              className="input-field"
              value={formData.bristolType}
              onChange={(e) =>
                handleInputChange("bristolType", parseInt(e.target.value) || 0)
              }
              required
            >
              <option value="">Select Bristol Type</option>
              <option value="1">
                Type 1 - Hard lumps (Severe constipation)
              </option>
              <option value="2">
                Type 2 - Lumpy sausage (Mild constipation)
              </option>
              <option value="3">Type 3 - Cracked sausage (Normal)</option>
              <option value="4">Type 4 - Smooth sausage (Ideal)</option>
              <option value="5">Type 5 - Soft blobs (Lacking fiber)</option>
              <option value="6">Type 6 - Fluffy pieces (Mild diarrhea)</option>
              <option value="7">Type 7 - Watery (Severe diarrhea)</option>
            </select>
          </div>

          <div className="grid md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Volume
              </label>
              <select
                className="input-field"
                value={formData.volume}
                onChange={(e) => handleInputChange("volume", e.target.value)}
              >
                <option value="">Select volume</option>
                <option value="Small">Small</option>
                <option value="Medium">Medium</option>
                <option value="Large">Large</option>
                <option value="Massive">Massive</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Color
              </label>
              <select
                className="input-field"
                value={formData.color}
                onChange={(e) => handleInputChange("color", e.target.value)}
              >
                <option value="">Select color</option>
                <option value="Brown">Brown</option>
                <option value="Dark Brown">Dark Brown</option>
                <option value="Light Brown">Light Brown</option>
                <option value="Yellow">Yellow</option>
                <option value="Green">Green</option>
                <option value="Red">Red</option>
                <option value="Black">Black</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Notes (Optional)
            </label>
            <textarea
              className="input-field"
              rows={3}
              placeholder="Any additional observations..."
              value={formData.notes}
              onChange={(e) => handleInputChange("notes", e.target.value)}
            ></textarea>
          </div>

          <button
            type="submit"
            className="btn-primary w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Saving..." : "Save Entry"}
          </button>
        </form>
      </div>

      {/* Recent Entries List */}
      <div className="mt-8">
        <h2 className="text-2xl font-bold mb-4">Recent Entries</h2>
        {isLoadingEntries ? (
          <div className="card text-center py-8">
            <p className="text-gray-600">Loading entries...</p>
          </div>
        ) : entries.length === 0 ? (
          <div className="card text-center py-8">
            <p className="text-gray-600 flex items-center justify-center gap-2">
              No entries yet. Create your first one above! <Logo size={24} />
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {entries.map((entry) => (
              <div key={entry.id} className="card">
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <h3 className="font-semibold text-lg">
                      Bristol Type {entry.bristolType}
                    </h3>
                    <p className="text-sm text-gray-600">
                      {getBristolTypeDescription(entry.bristolType)}
                    </p>
                  </div>
                  <span className="text-sm text-gray-500">
                    {new Date(entry.createdAt).toLocaleDateString()} at{" "}
                    {new Date(entry.createdAt).toLocaleTimeString([], {
                      hour: "2-digit",
                      minute: "2-digit"
                    })}
                  </span>
                </div>

                <div className="grid grid-cols-2 gap-4 mb-3 text-sm">
                  {entry.volume && (
                    <div>
                      <span className="font-medium">Volume:</span>{" "}
                      {entry.volume}
                    </div>
                  )}
                  {entry.color && (
                    <div>
                      <span className="font-medium">Color:</span> {entry.color}
                    </div>
                  )}
                </div>

                {entry.notes && (
                  <div className="mt-3 p-3 bg-gray-50 rounded">
                    <p className="text-sm text-gray-700">{entry.notes}</p>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
