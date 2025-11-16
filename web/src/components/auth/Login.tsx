import type React from "react";
import { LogIn } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import TopAppBar from "../TopAppBar";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  // TODO: backend integration
  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (!username.trim() || !password.trim()) {
      setError("Please fill in all fields");
      return;
    }

    // Placeholder: simulate login
    console.log("Login attempt:", { username, password });
    
    // Simulate successful login
    setTimeout(() => {
      // In real implementation, store token/session here
      console.log("Login successful (placeholder)");
      navigate("/");
    }, 500);
  }

  return (
    <div>
      <TopAppBar onSearch={(text) => console.log(text)} />
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Login
        </header>
        <div className="flex flex-col items-center justify-center min-h-[calc(100vh-64px)] p-6">
          <div className="w-full max-w-md space-y-6">
            <div className="flex flex-col items-center gap-4">
              <div className="bg-indigo-100 text-indigo-600 p-4 rounded-full">
                <LogIn className="w-10 h-10" />
              </div>
              <h1 className="text-3xl font-bold text-gray-900">Welcome back</h1>
              <p className="text-gray-500 text-center">
                Sign in to your account to continue
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="alert alert-error">
                  <span>{error}</span>
                </div>
              )}

              <div className="form-control">
                <label className="label" htmlFor="username">
                  <span className="label-text">Username</span>
                </label>
                <input
                  id="username"
                  type="text"
                  placeholder="Enter your username"
                  className="input input-bordered w-full"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>

              <div className="form-control">
                <label className="label" htmlFor="password">
                  <span className="label-text">Password</span>
                </label>
                <input
                  id="password"
                  type="password"
                  placeholder="Enter your password"
                  className="input input-bordered w-full"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>

              <div className="form-control mt-6">
                <button type="submit" className="btn btn-primary w-full">
                  Sign in
                </button>
              </div>
            </form>

            <div className="text-center">
              <p className="text-sm text-gray-600">
                Don't have an account?{" "}
                <Link to="/register" className="link link-primary">
                  Register here
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

