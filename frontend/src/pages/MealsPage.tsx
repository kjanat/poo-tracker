import { useState, useEffect } from "react";
import { useAuthStore } from "../stores/authStore";

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
  mealTime: string;
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
          setMeals(mealsData);
        }
      } catch (err) {
        console.error("Failed to fetch meals:", err);
      } finally {
        setLoadingMeals(false);
      }
    };

    fetchMeals();
  }, [token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      if (!token) {
        throw new Error("Not authenticated");
      }

      const mealData = {
        ...formData,
        mealTime: new Date().toISOString(), // Current time for now
        spicyLevel: formData.spicyLevel || null,
        // Convert empty strings to undefined for optional fields
        category: formData.category || undefined,
        description: formData.description || undefined,
        cuisine: formData.cuisine || undefined,
        notes: formData.notes || undefined
      };

      const response = await fetch(`${API_BASE_URL}/api/meals`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(mealData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to save meal");
      }

      const newMeal = await response.json();
      setMeals((prev) => [newMeal, ...prev]); // Add new meal to top of list
      setSuccess("🎉 Meal saved successfully!");

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
    } catch (err: any) {
      setError(err.message || "Failed to save meal");
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

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">🍽️ Meals</h1>

      <div className="card mb-6">
        <h3 className="text-lg font-semibold mb-4">Log New Meal</h3>

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
            ></textarea>
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
            ></textarea>
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

          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? "Saving..." : "Save Meal"}
          </button>
        </form>
      </div>

      {/* Existing Meals List */}
      <div className="card">
        <h3 className="text-lg font-semibold mb-4">Your Recent Meals</h3>

        {loadingMeals ? (
          <div className="text-center py-4">Loading meals...</div>
        ) : meals.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <p>No meals logged yet. Add your first meal above!</p>
          </div>
        ) : (
          <div className="space-y-4">
            {meals.map((meal) => (
              <div key={meal.id} className="border rounded-lg p-4 bg-gray-50">
                <div className="flex justify-between items-start mb-2">
                  <h4 className="font-semibold text-lg">{meal.name}</h4>
                  <div className="text-sm text-gray-500">
                    {new Date(meal.mealTime).toLocaleDateString()}{" "}
                    {new Date(meal.mealTime).toLocaleTimeString()}
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
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
