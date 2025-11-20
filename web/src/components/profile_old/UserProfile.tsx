import { UserX } from "lucide-react";
import { useState } from "react";
import { useLoaderData } from "react-router";
import { type Post, samplePosts } from "../feed/feedData";
import PostItem from "../feed/PostItem";

export function getUserInitials(username: string): string {
  if (!username) return "";
  return username.replace(/^@/, "").slice(0, 2).toUpperCase();
}
export default function UserProfile() {
  // TODO: backend req
  const [isFollowed, setIsFollowed] = useState(false);
  const user = useLoaderData();

  if (user === undefined) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen text-center p-6">
        <div className="flex flex-col items-center gap-4 max-w-md">
          <div className="bg-red-100 text-red-600 p-4 rounded-full">
            <UserX className="w-10 h-10" />
          </div>
          <h1 className="text-3xl font-bold">
            Sorry, that user does not exist
          </h1>
          <p className="text-gray-500">
            The user you’re looking for might have changed their username,
            deleted their account, or never existed at all.
          </p>
          <a
            href="/"
            className="mt-4 inline-block px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors"
          >
            Go back home
          </a>
        </div>
      </div>
    );
  }
  const userPosts = samplePosts.filter((p: Post) => p.author === user.username);
  // end of backend calls

  const initials = getUserInitials(user.username);

  //TODO backend for follow/unfollow V
  return (
    <div>
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Profile
        </header>
        <div className="p-6">
          <div className="flex flex-col items-center space-y-4">
            <div className="avatar avatar-placeholder">
              <div className="bg-neutral text-neutral-content w-24 rounded-full">
                <span className="text-3xl">{initials}</span>
              </div>
            </div>
            <div className="text-center">
              <div className="text-gray-600">{user.username}</div>
              <div className="text-2xl font-semibold text-gray-700">
                {user.bio}
              </div>
              <button
                type="button"
                onClick={() => setIsFollowed(!isFollowed)}
                className={`btn ${isFollowed ? "btn-outline btn-secondary" : "btn-primary"}`}
              >
                {isFollowed ? "Unfollow" : "Follow"}
              </button>
            </div>
          </div>
        </div>
        <div className="border-t border-gray-300" />
        {userPosts.map((post) => (
          <PostItem key={post.id} post={post} />
        ))}
      </div>
    </div>
  );
}
