import type React from "react";
import { UserPlus } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import TopAppBar from "../TopAppBar";

export default function Register() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  // TODO: backend integration
  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (!username.trim() || !password.trim() || !confirmPassword.trim()) {
      setError("Please fill in all fields");
      return;
    }

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (password.length < 6) {
      setError("Password must be at least 6 characters");
      return;
    }

    // Placeholder: simulate registration
    console.log("Registration attempt:", { username, password });
    
    // Simulate successful registration
    setTimeout(() => {
      // In real implementation, create user account here
      console.log("Registration successful (placeholder)");
      navigate("/login");
    }, 500);
  }

  return (
    <div>
      <TopAppBar onSearch={(text) => console.log(text)} />
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Register
        </header>
        <div className="flex flex-col items-center justify-center min-h-[calc(100vh-64px)] p-6">
          <div className="w-full max-w-md space-y-6">
            <div className="flex flex-col items-center gap-4">
              <div className="bg-green-100 text-green-600 p-4 rounded-full">
                <UserPlus className="w-10 h-10" />
              </div>
              <h1 className="text-3xl font-bold text-gray-900">Create account</h1>
              <p className="text-gray-500 text-center">
                Sign up to join the community
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="alert alert-error">
                  <span>{error}</span>
                </div>
              )}

              <div className="form-control">
                <label className="label" htmlFor="reg-username">
                  <span className="label-text">Username</span>
                </label>
                <input
                  id="reg-username"
                  type="text"
                  placeholder="Choose a username"
                  className="input input-bordered w-full"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>

              <div className="form-control">
                <label className="label" htmlFor="reg-password">
                  <span className="label-text">Password</span>
                </label>
                <input
                  id="reg-password"
                  type="password"
                  placeholder="Create a password"
                  className="input input-bordered w-full"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>

              <div className="form-control">
                <label className="label" htmlFor="confirm-password">
                  <span className="label-text">Confirm Password</span>
                </label>
                <input
                  id="confirm-password"
                  type="password"
                  placeholder="Confirm your password"
                  className="input input-bordered w-full"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                />
              </div>

              <div className="form-control mt-6">
                <button type="submit" className="btn btn-primary w-full">
                  Create account
                </button>
              </div>
            </form>

            <div className="text-center">
              <p className="text-sm text-gray-600">
                Already have an account?{" "}
                <Link to="/login" className="link link-primary">
                  Sign in here
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

