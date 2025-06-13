import { Link } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { BarChart3, Camera, Utensils, Brain } from 'lucide-react'

export function HomePage() {
  const { isAuthenticated } = useAuthStore()

  return (
    <div className="max-w-6xl mx-auto">
      {/* Hero Section */}
      <div className="text-center py-16">
        <h1 className="text-6xl font-bold mb-4">
          <span className="text-6xl mr-4">💩</span>
          Poo Tracker
        </h1>
        <p className="text-2xl text-gray-600 mb-8 max-w-3xl mx-auto">
          The brutally honest app that lets you log every majestic turd, rabbit pellet, 
          and volcanic diarrhea eruption without shame or censorship.
        </p>
        
        {!isAuthenticated ? (
          <div className="space-x-4">
            <Link to="/login" className="btn-primary text-lg px-8 py-3">
              Get Started
            </Link>
            <a 
              href="#features" 
              className="btn-secondary text-lg px-8 py-3"
            >
              Learn More
            </a>
          </div>
        ) : (
          <div className="space-x-4">
            <Link to="/dashboard" className="btn-primary text-lg px-8 py-3">
              Go to Dashboard
            </Link>
            <Link to="/new-entry" className="btn-secondary text-lg px-8 py-3">
              Log New Poop
            </Link>
          </div>
        )}
      </div>

      {/* Features Section */}
      <section id="features" className="py-16">
        <h2 className="text-4xl font-bold text-center mb-12">
          What's the fucking point?
        </h2>
        
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          <div className="card text-center">
            <BarChart3 className="mx-auto mb-4 text-poo-brown-600" size={48} />
            <h3 className="text-xl font-semibold mb-2">Track Your Shits</h3>
            <p className="text-gray-600">
              Log every glorious bowel movement with Bristol scores, timing, and satisfaction ratings.
            </p>
          </div>
          
          <div className="card text-center">
            <Camera className="mx-auto mb-4 text-poo-brown-600" size={48} />
            <h3 className="text-xl font-semibold mb-2">Photo Evidence</h3>
            <p className="text-gray-600">
              Snap a pic for science (or just to traumatize your friends). Visual documentation for the brave.
            </p>
          </div>
          
          <div className="card text-center">
            <Utensils className="mx-auto mb-4 text-poo-brown-600" size={48} />
            <h3 className="text-xl font-semibold mb-2">Meal Correlation</h3>
            <p className="text-gray-600">
              Record what you eat and discover if Taco Tuesday really is a war crime against your colon.
            </p>
          </div>
          
          <div className="card text-center">
            <Brain className="mx-auto mb-4 text-poo-brown-600" size={48} />
            <h3 className="text-xl font-semibold mb-2">AI Analysis</h3>
            <p className="text-gray-600">
              Our heartless machine learning model correlates patterns without judgment. Science is beautiful.
            </p>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-16 bg-white rounded-xl shadow-lg">
        <div className="px-8">
          <h2 className="text-4xl font-bold text-center mb-12">How it works</h2>
          
          <div className="grid md:grid-cols-4 gap-8 text-center">
            <div className="space-y-4">
              <div className="text-4xl">🍔</div>
              <h3 className="text-xl font-semibold">1. Eat</h3>
              <p className="text-gray-600">Something questionable.</p>
            </div>
            
            <div className="space-y-4">
              <div className="text-4xl">💩</div>
              <h3 className="text-xl font-semibold">2. Shit</h3>
              <p className="text-gray-600">Preferably in a toilet, but we're not here to kink-shame.</p>
            </div>
            
            <div className="space-y-4">
              <div className="text-4xl">📝</div>
              <h3 className="text-xl font-semibold">3. Log</h3>
              <p className="text-gray-600">Your experience in Poo Tracker with optional photos and ratings.</p>
            </div>
            
            <div className="space-y-4">
              <div className="text-4xl">🔄</div>
              <h3 className="text-xl font-semibold">4. Repeat</h3>
              <p className="text-gray-600">Until you realize you're lactose intolerant or have IBS.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Privacy Section */}
      <section className="py-16 text-center">
        <h2 className="text-3xl font-bold mb-6">Privacy, baby</h2>
        <p className="text-lg text-gray-600 max-w-2xl mx-auto">
          We encrypt your brown notes and hide them away. Nobody's reading your logs except you—and 
          whatever godforsaken AI wants to learn about the day-after effects of your sushi buffet.
        </p>
      </section>

      {/* Disclaimer */}
      <section className="py-8 text-center border-t border-gray-200">
        <p className="text-sm text-gray-500">
          <strong>Disclaimer:</strong> Not responsible for phone screen damage caused by ill-advised photo documentation. 
          Use with pride, shame, or scientific detachment. Up to you.
        </p>
      </section>
    </div>
  )
}
