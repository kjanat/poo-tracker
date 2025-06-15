import Logo from "../components/Logo";
import { getLogoProps } from "../utils/branding";

export function LoginPage() {
  return (
    <div className="max-w-md mx-auto">
      <div className="card">
        <div className="text-center mb-6">
          <Logo {...getLogoProps("login")} className="mx-auto mb-4" />
          <h1 className="text-2xl font-bold">Welcome Back</h1>
        </div>
        <p className="text-center text-gray-600 mb-8">
          Ready to log some legendary logs?
        </p>

        <form className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              className="input-field"
              placeholder="your@email.com"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              type="password"
              className="input-field"
              placeholder="Password"
              required
            />
          </div>

          <button type="submit" className="btn-primary w-full">
            Sign In
          </button>
        </form>

        <p className="text-center text-sm text-gray-600 mt-6">
          Don't have an account? We'll create one for you automatically on first
          login.
        </p>
      </div>
    </div>
  );
}
