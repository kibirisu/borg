import { NavLink } from "react-router";
import { useContext, type ReactNode } from "react";
import AppContext from "../../lib/state";

type SidebarItem = {
  label: string;
  to: string;
  icon: ReactNode;
};

const items: SidebarItem[] = [
  {
    label: "Home",
    to: "/",
    icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24">
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M3 12l9-9 9 9M4 10v10a1 1 0 0 0 1 1h5v-6h4v6h5a1 1 0 0 0 1-1V10"
        />
      </svg>
    ),
  },
  {
    label: "Explore",
    to: "/explore",
    icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24">
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M10 6.025A7.5 7.5 0 1 0 17.975 14H10V6.025Z"
        />
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M13.5 3V11h7.975"
        />
      </svg>
    ),
  },
  {
    label: "Notifications",
    to: "/notifications",
    icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24">
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M15 17h5l-1.4-1.4A2 2 0 0 1 18 14.2V11a6 6 0 1 0-12 0v3.2a2 2 0 0 1-.6 1.4L4 17h5"
        />
      </svg>
    ),
  },
  {
    label: "Likes",
    to: "/likes",
    icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24">
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.6l-1-1a5.5 5.5 0 0 0-7.8 7.8l1 1L12 21l7.8-7.8 1-1a5.5 5.5 0 0 0 0-7.8Z"
        />
      </svg>
    ),
  },
  {
    label: "Shared",
    to: "/shared",
    icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24">
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M4 12v7a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-7"
        />
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M16 6l-4-4-4 4M12 2v14"
        />
      </svg>
    ),
  },
  {
    label: "Profile",
    to: "/profile",
    icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24">
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M12 12a5 5 0 1 0-5-5 5 5 0 0 0 5 5Z"
        />
        <path
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M3 21a9 9 0 0 1 18 0"
        />
      </svg>
    ),
  },
  {
    label: "More",
    to: "/more",
    icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24">
        <circle cx="5" cy="12" r="1.5" fill="currentColor" />
        <circle cx="12" cy="12" r="1.5" fill="currentColor" />
        <circle cx="19" cy="12" r="1.5" fill="currentColor" />
      </svg>
    ),
  },
];

export default function Sidebar() {
  const appState = useContext(AppContext);
  const profileTarget = appState?.userId
    ? `/profile/${appState.userId}`
    : "/signin";

  return (
    <aside className="w-64 h-screen bg-white border-r px-4 py-6">
      <h1 className="text-xl font-bold mb-8">Borg</h1>

      <nav className="space-y-2">
        {items.map((item) => {
          const destination = item.to === "/profile" ? profileTarget : item.to;
          const end = item.to === "/profile";
          return (
            <NavLink
              key={item.label}
              to={destination}
              end={end}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2 rounded-lg font-medium
              ${
                  isActive
                    ? "bg-indigo-50 text-indigo-600"
                    : "text-gray-700 hover:bg-gray-100"
                }`
              }
            >
              {item.icon}
              <span>{item.label}</span>
            </NavLink>
          );
        })}
      </nav>

      <button className="mt-8 w-full bg-indigo-500 text-white py-2 rounded-full font-semibold">
        Post
      </button>
    </aside>
  );
}
