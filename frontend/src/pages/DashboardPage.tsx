export function DashboardPage() {
  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">💩 Dashboard</h1>
      
      <div className="grid md:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">Total Entries</h3>
          <p className="text-3xl font-bold text-poo-brown-600">23</p>
        </div>
        
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">Avg Bristol Score</h3>
          <p className="text-3xl font-bold text-poo-brown-600">4.2</p>
        </div>
        
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">This Week</h3>
          <p className="text-3xl font-bold text-poo-brown-600">7</p>
        </div>
      </div>
      
      <div className="grid md:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">Recent Entries</h3>
          <p className="text-gray-600">Your recent bowel movements will appear here...</p>
        </div>
        
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">Quick Stats</h3>
          <p className="text-gray-600">Analytics and trends will be displayed here...</p>
        </div>
      </div>
    </div>
  )
}
