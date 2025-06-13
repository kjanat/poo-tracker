export function MealsPage() {
  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">🍽️ Meals</h1>
      
      <div className="card mb-6">
        <h3 className="text-lg font-semibold mb-4">Log New Meal</h3>
        
        <form className="space-y-4">
          <div className="grid md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Meal Name
              </label>
              <input 
                type="text" 
                className="input-field" 
                placeholder="e.g., Spicy Tacos"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Category
              </label>
              <select className="input-field">
                <option value="">Select category</option>
                <option value="Breakfast">Breakfast</option>
                <option value="Lunch">Lunch</option>
                <option value="Dinner">Dinner</option>
                <option value="Snack">Snack</option>
              </select>
            </div>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description (Optional)
            </label>
            <textarea 
              className="input-field" 
              rows={2}
              placeholder="Describe the meal..."
            ></textarea>
          </div>
          
          <button type="submit" className="btn-primary">
            Save Meal
          </button>
        </form>
      </div>
      
      <div className="card">
        <h3 className="text-lg font-semibold mb-4">Recent Meals</h3>
        <p className="text-gray-600">Your meal history will appear here...</p>
      </div>
    </div>
  )
}
