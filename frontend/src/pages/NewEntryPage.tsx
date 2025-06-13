export function NewEntryPage() {
  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Log New Entry</h1>
      
      <div className="card">
        <p className="text-center text-gray-600 mb-8">
          Time to document another masterpiece! 💩
        </p>
        
        <form className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Bristol Stool Chart Type (1-7)
            </label>
            <select className="input-field">
              <option value="">Select Bristol Type</option>
              <option value="1">Type 1 - Hard lumps (Severe constipation)</option>
              <option value="2">Type 2 - Lumpy sausage (Mild constipation)</option>
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
              <select className="input-field">
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
              <select className="input-field">
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
            ></textarea>
          </div>
          
          <button type="submit" className="btn-primary w-full">
            Save Entry
          </button>
        </form>
      </div>
    </div>
  )
}
