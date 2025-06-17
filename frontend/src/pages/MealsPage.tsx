import { useState, useEffect } from "react";
import { useAuthStore } from "../stores/authStore";
import Logo from "../components/Logo";

interface Meal {
  id: string;
  name: string;
  category?: string;
  description?: string;
  cuisine?: string;
  spicyLevel?: number;
  fiberRich: boolean;
  dairy: boolean;
  gluten: boolean;
  notes?: string;
  photoUrl?: string;
  mealTime: string;
  createdAt: string;
  linkedEntries?: Entry[];
}

interface Entry {
  id: string;
  bristolType: number;
  volume?: string;
  color?: string;
  consistency?: string;
  notes?: string;
  createdAt: string;
}

interface MealFormData {
  name: string;
  category: string;
  description: string;
  cuisine: string;
  spicyLevel: number;
  fiberRich: boolean;
  dairy: boolean;
  gluten: boolean;
  notes: string;
  photoUrl?: string;
}

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3002";

export function MealsPage() {
  const [formData, setFormData] = useState<MealFormData>({
    name: "",
    category: "",
    description: "",
    cuisine: "",
    spicyLevel: 1,
    fiberRich: false,
    dairy: false,
    gluten: false,
    notes: ""
  });
  const [meals, setMeals] = useState<Meal[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMeals, setLoadingMeals] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [editingMeal, setEditingMeal] = useState<Meal | null>(null);
  const [linkingMeal, setLinkingMeal] = useState<Meal | null>(null);
  const [availableEntries, setAvailableEntries] = useState<Entry[]>([]);
  const [linkedEntries, setLinkedEntries] = useState<Entry[]>([]);
  const [showLinkingModal, setShowLinkingModal] = useState(false);

  const token = useAuthStore((state) => state.token);

  // Fetch existing meals
  useEffect(() => {
    const fetchMeals = async () => {
      if (!token) return;

      try {
        const response = await fetch(`${API_BASE_URL}/api/meals`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });

        if (response.ok) {
          const mealsData = await response.json();

          // Fetch linked entries for each meal
          const mealsWithLinkedEntries = await Promise.all(
            mealsData.map(async (meal: Meal) => {
              try {
                const linkedResponse = await fetch(
                  `${API_BASE_URL}/api/meals/${meal.id}/entries`,
                  {
                    headers: { Authorization: `Bearer ${token}` }
                  }
                );

                if (linkedResponse.ok) {
                  const linkedData = await linkedResponse.json();
                  return { ...meal, linkedEntries: linkedData };
                } else {
                  return { ...meal, linkedEntries: [] };
                }
              } catch (err) {
                console.error(
                  `Failed to fetch linked entries for meal ${meal.id}:`,
                  err
                );
                return { ...meal, linkedEntries: [] };
              }
            })
          );

          setMeals(mealsWithLinkedEntries);
        }
      } catch (err) {
        console.error("Failed to fetch meals:", err);
      } finally {
        setLoadingMeals(false);
      }
    };

    fetchMeals();
  }, [token]);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      // Validate file type
      if (!file.type.startsWith("image/")) {
        setError("Please select a valid image file");
        return;
      }

      // Validate file size (5MB limit)
      if (file.size > 5 * 1024 * 1024) {
        setError("Image size must be less than 5MB");
        return;
      }

      setSelectedImage(file);

      // Create preview
      const reader = new FileReader();
      reader.onload = (e) => {
        setImagePreview(e.target?.result as string);
      };
      reader.readAsDataURL(file);

      // Clear any previous errors
      setError("");
    }
  };

  const removeImage = () => {
    setSelectedImage(null);
    setImagePreview(null);
  };

  const uploadImage = async (): Promise<string | null> => {
    if (!selectedImage) return null;

    const uploadFormData = new FormData();
    uploadFormData.append("image", selectedImage);

    try {
      const response = await fetch(`${API_BASE_URL}/api/uploads`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`
        },
        body: uploadFormData
      });

      if (!response.ok) {
        throw new Error("Failed to upload image");
      }

      const data = await response.json();
      return data.imageUrl;
    } catch (error) {
      console.error("Error uploading image:", error);
      throw error;
    }
  };

  const startEdit = (meal: Meal) => {
    setEditingMeal(meal);
    setFormData({
      name: meal.name,
      category: meal.category || "",
      description: meal.description || "",
      cuisine: meal.cuisine || "",
      spicyLevel: meal.spicyLevel || 1,
      fiberRich: meal.fiberRich,
      dairy: meal.dairy,
      gluten: meal.gluten,
      notes: meal.notes || "",
      photoUrl: meal.photoUrl
    });
    // Reset image states when editing
    setSelectedImage(null);
    setImagePreview(null);
    // Scroll to form
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const cancelEdit = () => {
    setEditingMeal(null);
    setFormData({
      name: "",
      category: "",
      description: "",
      cuisine: "",
      spicyLevel: 1,
      fiberRich: false,
      dairy: false,
      gluten: false,
      notes: ""
    });
    setSelectedImage(null);
    setImagePreview(null);
    setError("");
    setSuccess("");
  };

  const deleteMeal = async (mealId: string) => {
    if (!confirm("Are you sure you want to delete this meal?")) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/meals/${mealId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error("Failed to delete meal");
      }

      // Remove from local state
      setMeals((prev) => prev.filter((meal) => meal.id !== mealId));
      setSuccess("Meal deleted successfully!");
    } catch (error) {
      console.error("Error deleting meal:", error);
      setError(
        error instanceof Error ? error.message : "Failed to delete meal"
      );
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      if (!token) {
        throw new Error("Not authenticated");
      }

      // Upload image if selected
      let photoUrl = formData.photoUrl;
      if (selectedImage) {
        const uploadedUrl = await uploadImage();
        if (uploadedUrl) {
          photoUrl = uploadedUrl;
        }
      }

      const mealData = {
        ...formData,
        photoUrl: photoUrl || undefined,
        mealTime: editingMeal ? editingMeal.mealTime : new Date().toISOString(),
        spicyLevel: formData.spicyLevel || null,
        // Convert empty strings to undefined for optional fields
        category: formData.category || undefined,
        description: formData.description || undefined,
        cuisine: formData.cuisine || undefined,
        notes: formData.notes || undefined
      };

      const isEditing = editingMeal !== null;
      const url = isEditing
        ? `${API_BASE_URL}/api/meals/${editingMeal.id}`
        : `${API_BASE_URL}/api/meals`;
      const method = isEditing ? "PUT" : "POST";

      const response = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(mealData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(
          errorData.error || `Failed to ${isEditing ? "update" : "save"} meal`
        );
      }

      const savedMeal = await response.json();

      if (isEditing) {
        // Update the meal in the list
        setMeals((prev) =>
          prev.map((meal) => (meal.id === editingMeal.id ? savedMeal : meal))
        );
        setSuccess("🎉 Meal updated successfully!");
        setEditingMeal(null);
      } else {
        // Add new meal to top of list
        setMeals((prev) => [savedMeal, ...prev]);
        setSuccess("🎉 Meal saved successfully!");
      }

      // Reset form
      setFormData({
        name: "",
        category: "",
        description: "",
        cuisine: "",
        spicyLevel: 1,
        fiberRich: false,
        dairy: false,
        gluten: false,
        notes: ""
      });
      setSelectedImage(null);
      setImagePreview(null);
    } catch (err: any) {
      setError(
        err.message || `Failed to ${editingMeal ? "update" : "save"} meal`
      );
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value, type } = e.target;

    if (type === "checkbox") {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData((prev) => ({ ...prev, [name]: checked }));
    } else if (type === "number") {
      setFormData((prev) => ({ ...prev, [name]: parseInt(value) || 1 }));
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
    }
  };

  // Functions for linking entries to meals
  const openLinkingModal = async (meal: Meal) => {
    setLinkingMeal(meal);
    setShowLinkingModal(true);

    try {
      // Fetch available entries
      const entriesResponse = await fetch(`${API_BASE_URL}/api/entries`, {
        headers: { Authorization: `Bearer ${token}` }
      });

      if (entriesResponse.ok) {
        const entriesData = await entriesResponse.json();
        setAvailableEntries(entriesData.entries || []);
      }

      // Fetch already linked entries
      const linkedResponse = await fetch(
        `${API_BASE_URL}/api/meals/${meal.id}/entries`,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );

      if (linkedResponse.ok) {
        const linkedData = await linkedResponse.json();
        setLinkedEntries(linkedData);
      }
    } catch (err) {
      console.error("Failed to fetch entries:", err);
    }
  };

  const linkEntryToMeal = async (entryId: string) => {
    if (!linkingMeal) return;

    try {
      const response = await fetch(
        `${API_BASE_URL}/api/meals/${linkingMeal.id}/link-entry`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
          },
          body: JSON.stringify({ entryId })
        }
      );

      if (response.ok) {
        // Refresh linked entries in the modal
        const linkedResponse = await fetch(
          `${API_BASE_URL}/api/meals/${linkingMeal.id}/entries`,
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        );

        if (linkedResponse.ok) {
          const linkedData = await linkedResponse.json();
          setLinkedEntries(linkedData);

          // Also update the main meals list to reflect the change
          setMeals((prevMeals) =>
            prevMeals.map((meal) =>
              meal.id === linkingMeal.id
                ? { ...meal, linkedEntries: linkedData }
                : meal
            )
          );
        }
        setSuccess("Entry linked successfully!");
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Failed to link entry");
      }
    } catch (err) {
      setError("Failed to link entry");
    }
  };

  const unlinkEntryFromMeal = async (entryId: string) => {
    if (!linkingMeal) return;

    try {
      const response = await fetch(
        `${API_BASE_URL}/api/meals/${linkingMeal.id}/unlink-entry`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
          },
          body: JSON.stringify({ entryId })
        }
      );

      if (response.ok) {
        // Refresh linked entries in the modal
        const linkedResponse = await fetch(
          `${API_BASE_URL}/api/meals/${linkingMeal.id}/entries`,
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        );

        if (linkedResponse.ok) {
          const linkedData = await linkedResponse.json();
          setLinkedEntries(linkedData);

          // Also update the main meals list to reflect the change
          setMeals((prevMeals) =>
            prevMeals.map((meal) =>
              meal.id === linkingMeal.id
                ? { ...meal, linkedEntries: linkedData }
                : meal
            )
          );
        }
        setSuccess("Entry unlinked successfully!");
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Failed to unlink entry");
      }
    } catch (err) {
      setError("Failed to unlink entry");
    }
  };

  // Function to unlink entry directly from the main meal view
  const unlinkEntryFromMainView = async (mealId: string, entryId: string) => {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/meals/${mealId}/unlink-entry`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
          },
          body: JSON.stringify({ entryId })
        }
      );

      if (response.ok) {
        // Refresh linked entries for this specific meal
        const linkedResponse = await fetch(
          `${API_BASE_URL}/api/meals/${mealId}/entries`,
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        );

        if (linkedResponse.ok) {
          const linkedData = await linkedResponse.json();

          // Update the main meals list to reflect the change
          setMeals((prevMeals) =>
            prevMeals.map((meal) =>
              meal.id === mealId ? { ...meal, linkedEntries: linkedData } : meal
            )
          );
        }
        setSuccess("Entry unlinked successfully!");

        // Clear success message after a few seconds
        setTimeout(() => setSuccess(""), 3000);
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Failed to unlink entry");
      }
    } catch (err) {
      setError("Failed to unlink entry");
    }
  };

  const formatBristolType = (type: number) => {
    const types = [
      "Separate hard lumps",
      "Lumpy sausage",
      "Cracked sausage",
      "Smooth snake",
      "Soft blobs",
      "Mushy stool",
      "Liquid stool"
    ];
    return `Type ${type} - ${types[type - 1]}`;
  };

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8 flex items-center gap-2">
        <Logo size={32} /> Meals
      </h1>

      <div className="card mb-6">
        <h3 className="text-lg font-semibold mb-4">
          {editingMeal ? "Edit Meal" : "Log New Meal"}
        </h3>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {success && (
          <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
            {success}
          </div>
        )}

        <form className="space-y-4" onSubmit={handleSubmit}>
          <div className="grid md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Meal Name *
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleInputChange}
                className="input-field"
                placeholder="e.g., Spicy Tacos"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Category
              </label>
              <select
                name="category"
                value={formData.category}
                onChange={handleInputChange}
                className="input-field"
              >
                <option value="">Select category</option>
                <option value="Breakfast">Breakfast</option>
                <option value="Lunch">Lunch</option>
                <option value="Dinner">Dinner</option>
                <option value="Snack">Snack</option>
              </select>
            </div>
          </div>

          <div className="grid md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Cuisine
              </label>
              <input
                type="text"
                name="cuisine"
                value={formData.cuisine}
                onChange={handleInputChange}
                className="input-field"
                placeholder="e.g., Mexican, Italian, Thai"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Spicy Level (1-10)
              </label>
              <input
                type="number"
                name="spicyLevel"
                value={formData.spicyLevel}
                onChange={handleInputChange}
                className="input-field"
                min="1"
                max="10"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleInputChange}
              className="input-field"
              rows={2}
              placeholder="Describe the meal..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Additional Notes
            </label>
            <textarea
              name="notes"
              value={formData.notes}
              onChange={handleInputChange}
              className="input-field"
              rows={2}
              placeholder="Any additional notes..."
            />
          </div>

          {/* Image Upload Section */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Photo (Optional)
            </label>
            <div className="space-y-3">
              <input
                type="file"
                accept="image/*"
                onChange={handleImageChange}
                className="input-field"
              />

              {imagePreview && (
                <div className="relative">
                  <img
                    src={imagePreview}
                    alt="Preview"
                    className="max-w-xs max-h-48 object-cover rounded border"
                  />
                  <button
                    type="button"
                    onClick={removeImage}
                    className="absolute top-2 right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-sm hover:bg-red-600"
                  >
                    ×
                  </button>
                </div>
              )}

              {editingMeal?.photoUrl && !imagePreview && (
                <div className="text-sm text-gray-600">
                  Current photo:
                  <img
                    src={`${API_BASE_URL}${editingMeal.photoUrl}`}
                    alt="Current meal photo"
                    className="mt-2 max-w-xs max-h-48 object-cover rounded border"
                  />
                  <p className="mt-1 text-xs">
                    Upload a new image to replace it
                  </p>
                </div>
              )}

              <p className="text-xs text-gray-500">
                Max file size: 5MB. Accepted formats: JPG, PNG, GIF, WebP
              </p>
            </div>
          </div>

          <div className="grid md:grid-cols-3 gap-4">
            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                name="fiberRich"
                checked={formData.fiberRich}
                onChange={handleInputChange}
                className="rounded"
              />
              <span className="text-sm">High Fiber</span>
            </label>

            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                name="dairy"
                checked={formData.dairy}
                onChange={handleInputChange}
                className="rounded"
              />
              <span className="text-sm">Contains Dairy</span>
            </label>

            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                name="gluten"
                checked={formData.gluten}
                onChange={handleInputChange}
                className="rounded"
              />
              <span className="text-sm">Contains Gluten</span>
            </label>
          </div>

          <div className="flex gap-3">
            <button
              type="submit"
              className="btn-primary flex-1"
              disabled={loading}
            >
              {loading
                ? editingMeal
                  ? "Updating..."
                  : "Saving..."
                : editingMeal
                  ? "Update Meal"
                  : "Save Meal"}
            </button>

            {editingMeal && (
              <button
                type="button"
                onClick={cancelEdit}
                className="btn-secondary"
                disabled={loading}
              >
                Cancel
              </button>
            )}
          </div>
        </form>
      </div>

      {/* Existing Meals List */}
      <div className="card">
        <h3 className="text-lg font-semibold mb-4">Your Recent Meals</h3>

        {loadingMeals ? (
          <div className="text-center py-4">Loading meals...</div>
        ) : meals.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <p className="flex items-center justify-center gap-2">
              No meals logged yet. Add your first meal above! <Logo size={24} />
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {meals.map((meal) => (
              <div key={meal.id} className="border rounded-lg p-4 bg-gray-50">
                <div className="flex justify-between items-start mb-2">
                  <h4 className="font-semibold text-lg">{meal.name}</h4>
                  <div className="flex items-center gap-2">
                    <div className="text-sm text-gray-500">
                      {new Date(meal.mealTime).toLocaleDateString()}{" "}
                      {new Date(meal.mealTime).toLocaleTimeString()}
                    </div>
                    <div className="flex gap-1">
                      <button
                        onClick={() => startEdit(meal)}
                        className="text-blue-600 hover:text-blue-800 text-sm font-medium"
                        disabled={loading}
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => openLinkingModal(meal)}
                        className="text-green-600 hover:text-green-800 text-sm font-medium"
                        disabled={loading}
                      >
                        Link Entries
                      </button>
                      <button
                        onClick={() => deleteMeal(meal.id)}
                        className="text-red-600 hover:text-red-800 text-sm font-medium"
                        disabled={loading}
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>

                <div className="grid md:grid-cols-2 gap-4 text-sm">
                  <div>
                    {meal.category && (
                      <p>
                        <strong>Category:</strong> {meal.category}
                      </p>
                    )}
                    {meal.cuisine && (
                      <p>
                        <strong>Cuisine:</strong> {meal.cuisine}
                      </p>
                    )}
                    {meal.spicyLevel && (
                      <p>
                        <strong>Spicy Level:</strong> {meal.spicyLevel}/10
                      </p>
                    )}
                  </div>
                  <div>
                    {meal.description && (
                      <p>
                        <strong>Description:</strong> {meal.description}
                      </p>
                    )}
                    {meal.notes && (
                      <p>
                        <strong>Notes:</strong> {meal.notes}
                      </p>
                    )}
                  </div>
                </div>

                {meal.photoUrl && (
                  <div className="mt-3">
                    <img
                      src={`${API_BASE_URL}${meal.photoUrl}`}
                      alt="Meal photo"
                      className="max-w-xs max-h-48 object-cover rounded border"
                    />
                  </div>
                )}

                {(meal.fiberRich || meal.dairy || meal.gluten) && (
                  <div className="mt-2 flex flex-wrap gap-2">
                    {meal.fiberRich && (
                      <span className="inline-block bg-green-100 text-green-800 text-xs px-2 py-1 rounded">
                        High Fiber
                      </span>
                    )}
                    {meal.dairy && (
                      <span className="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded">
                        Dairy
                      </span>
                    )}
                    {meal.gluten && (
                      <span className="inline-block bg-yellow-100 text-yellow-800 text-xs px-2 py-1 rounded">
                        Gluten
                      </span>
                    )}
                  </div>
                )}

                {/* Linked Entries Section */}
                <div className="mt-4">
                  <h5 className="font-semibold text-md mb-2">Linked Entries</h5>

                  {!meal.linkedEntries || meal.linkedEntries.length === 0 ? (
                    <p className="text-sm text-gray-500">
                      No entries linked to this meal.
                    </p>
                  ) : (
                    <div className="space-y-2">
                      {meal.linkedEntries.map((entry) => (
                        <div
                          key={entry.id}
                          className="p-3 rounded-lg bg-gray-100 flex justify-between items-center"
                        >
                          <div className="text-sm">
                            <p>
                              <strong>Bristol Type:</strong>{" "}
                              {formatBristolType(entry.bristolType)}
                            </p>
                            {entry.notes && (
                              <p>
                                <strong>Notes:</strong> {entry.notes}
                              </p>
                            )}
                          </div>

                          <button
                            onClick={() =>
                              unlinkEntryFromMainView(meal.id, entry.id)
                            }
                            className="text-red-600 hover:text-red-800 text-sm font-medium"
                            disabled={loading}
                          >
                            Unlink
                          </button>
                        </div>
                      ))}
                    </div>
                  )}

                  <button
                    onClick={() => openLinkingModal(meal)}
                    className="mt-2 text-blue-600 hover:text-blue-800 text-sm font-medium"
                    disabled={loading}
                  >
                    Link Entries
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Linking Modal */}
      {showLinkingModal && linkingMeal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-lg p-6 max-w-2xl w-full mx-4 max-h-[80vh] overflow-y-auto">
            <h3 className="text-lg font-semibold mb-4">
              Link Entries to "{linkingMeal.name}"
            </h3>

            {error && (
              <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
                {error}
              </div>
            )}

            {success && (
              <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
                {success}
              </div>
            )}

            {/* Already Linked Entries */}
            <div className="mb-6">
              <h4 className="font-medium mb-2">
                Already Linked Entries ({linkedEntries.length})
              </h4>
              {linkedEntries.length === 0 ? (
                <p className="text-sm text-gray-500 bg-gray-50 p-3 rounded">
                  No entries linked to this meal yet.
                </p>
              ) : (
                <div className="space-y-2">
                  {linkedEntries.map((entry) => (
                    <div
                      key={entry.id}
                      className="p-3 rounded-lg bg-blue-50 border border-blue-200 flex justify-between items-center"
                    >
                      <div className="text-sm">
                        <p className="font-medium">
                          {formatBristolType(entry.bristolType)}
                        </p>
                        <p className="text-gray-600">
                          {new Date(entry.createdAt).toLocaleDateString()}{" "}
                          {new Date(entry.createdAt).toLocaleTimeString()}
                        </p>
                        {entry.volume && (
                          <p className="text-gray-600">
                            Volume: {entry.volume}
                          </p>
                        )}
                        {entry.notes && (
                          <p className="text-gray-600">Notes: {entry.notes}</p>
                        )}
                      </div>
                      <button
                        onClick={() => unlinkEntryFromMeal(entry.id)}
                        className="text-red-600 hover:text-red-800 text-sm font-medium px-2 py-1 rounded"
                        disabled={loading}
                      >
                        Unlink
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Available Entries to Link */}
            <div className="mb-6">
              <h4 className="font-medium mb-2">Available Entries to Link</h4>
              {availableEntries.length === 0 ? (
                <p className="text-sm text-gray-500 bg-gray-50 p-3 rounded">
                  No available entries to link. Create some stool entries first!
                </p>
              ) : (
                <div className="space-y-2 max-h-60 overflow-y-auto">
                  {availableEntries
                    .filter(
                      (entry) =>
                        !linkedEntries.some((linked) => linked.id === entry.id)
                    )
                    .map((entry) => (
                      <div
                        key={entry.id}
                        className="p-3 rounded-lg bg-green-50 border border-green-200 flex justify-between items-center"
                      >
                        <div className="text-sm">
                          <p className="font-medium">
                            {formatBristolType(entry.bristolType)}
                          </p>
                          <p className="text-gray-600">
                            {new Date(entry.createdAt).toLocaleDateString()}{" "}
                            {new Date(entry.createdAt).toLocaleTimeString()}
                          </p>
                          {entry.volume && (
                            <p className="text-gray-600">
                              Volume: {entry.volume}
                            </p>
                          )}
                          {entry.notes && (
                            <p className="text-gray-600">
                              Notes: {entry.notes}
                            </p>
                          )}
                        </div>
                        <button
                          onClick={() => linkEntryToMeal(entry.id)}
                          className="text-green-600 hover:text-green-800 text-sm font-medium px-2 py-1 rounded bg-green-100 hover:bg-green-200"
                          disabled={loading}
                        >
                          Link
                        </button>
                      </div>
                    ))}
                </div>
              )}
            </div>

            {/* Modal Footer */}
            <div className="flex justify-end gap-3 pt-4 border-t">
              <button
                onClick={() => {
                  setShowLinkingModal(false);
                  setLinkingMeal(null);
                  setError("");
                  setSuccess("");
                }}
                className="px-4 py-2 text-gray-700 bg-gray-200 rounded hover:bg-gray-300 transition-colors"
                disabled={loading}
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
