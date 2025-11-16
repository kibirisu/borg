import { Link } from 'react-router';
import { useId, useState } from 'react';

const Input = (onSearch) => {
  const handleKeyDown = (event) => {
    if (event.key === 'Enter') {
      onSearch(event.currentTarget.value);
    }
  };

  return (
    <input
      type="text"
      placeholder="Search..."
      className="input input-bordered w-40 md:w-64"
      onKeyDown={handleKeyDown}
    />
  );
};
export default function TopAppBar({ onSearch }: Props) {
  const [showSearch, setShowSearch] = useState(false);

  return (
    <div className="navbar bg-base-100 shadow-sm z-100 sticky top-0">
      <div className="navbar-start">
        <div className="dropdown">
          <button type="button" className="btn btn-ghost btn-circle">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              {' '}
              <title id={useId()}>Menu Hamburger</title>
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h7" />{' '}
            </svg>
          </button>
          <ul tabIndex="-1" className="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow">
            <li className="menu-disabled">
              <a href="/profile">Profile</a>
            </li>
            <li className="menu-disabled">
              <a href="/trending">Trending</a>
            </li>
            <li className="menu-disabled">
              <a href="/federation">Federation</a>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
            <li>
              <Link to="/register">Register</Link>
            </li>
          </ul>
        </div>
      </div>
      <div className="navbar-center">
        <a href="/" className="btn btn-ghost text-xl">
          Name of our app :c
        </a>
      </div>
      <div className="navbar-end">
        <button type="button" className="btn btn-ghost btn-circle" onClick={() => setShowSearch(!showSearch)}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            {' '}
            <title id={useId()}>Search</title>
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />{' '}
          </svg>
        </button>
        {showSearch && <div className="form-control transition-all duration-300">{Input(onSearch)}</div>}
        <button type="button" className="btn btn-ghost btn-circle">
          <div className="indicator">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              {' '}
              <title id={useId()}>Notifications</title>
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
              />{' '}
            </svg>
            <span className="badge badge-xs badge-primary indicator-item"></span>
          </div>
        </button>
      </div>
    </div>
  );
}
