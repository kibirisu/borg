import { useQuery } from "@tanstack/react-query";
import { useLoaderData } from "react-router";
import type { components } from "../../lib/api/v1";
import type { Client } from "../../lib/client";
import { type Post } from "./feedData";
import NewPostBox from "./NewPostBox";
import PostItem from "./PostItem";

export const loader = (client: Client) => async () => {
  // Najprostszy loader - tylko zwraca opts, bez prefetch
  // Komponent użyje useQuery do pobrania danych
  try {
    const opts = client.$api.queryOptions('get', '/api/posts', {});
    if (!opts) {
      throw new Error('queryOptions returned undefined');
    }
    return { opts };
  } catch (error) {
    console.error('Loader error:', error);
    // Jeśli queryOptions nie działa, zwróć pusty obiekt z undefined opts
    // Komponent to obsłuży
    return { opts: undefined as any };
  }
};

export default function MainFeed() {
  const loaderData = useLoaderData() as Awaited<ReturnType<ReturnType<typeof loader>>> | undefined;
  console.log('MainFeed component - loaderData:', loaderData);
  
  if (!loaderData || !loaderData.opts) {
    return (
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Home
        </header>
        <div className="p-4 text-center text-red-600">
          Error: Loader data is missing. Please check if the backend endpoint /api/posts is implemented.
        </div>
      </div>
    );
  }
  
  const { opts } = loaderData;
  const { data: apiPosts, isPending, error } = useQuery(opts);

  // Mapowanie danych z API do typu Post używanego przez PostItem
  const posts: Post[] = apiPosts && Array.isArray(apiPosts)
    ? (apiPosts as components['schemas']['Post'][]).map((p) => {
        // Obsługa createdAt - zawsze string w API
        const createdAtStr = typeof p.createdAt === 'string' 
          ? p.createdAt 
          : new Date(p.createdAt).toISOString();

        return {
          id: String(p.id),
          author: p.username ? `@${p.username}` : `@user${p.userID}`,
          content: p.content,
          createdAt: createdAtStr,
          likes: p.likeCount,
          replies: p.commentCount,
          reposts: p.shareCount,
        };
      })
    : [];

  function addPost(content: string) {
    // TODO: Zaimplementować tworzenie posta przez API
    // Na razie funkcja pozostaje pusta, można później dodać mutację
    console.log('New post:', content);
  }

  if (isPending) {
    return (
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Home
        </header>
        <div className="p-4 text-center text-gray-500">Loading posts...</div>
      </div>
    );
  }

  if (error) {
    const errorMessage = (error as Error)?.message || String(error);
    return (
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Home
        </header>
        <div className="p-4 text-center text-red-600">
          Error loading posts: {errorMessage}
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
        Home
      </header>
      <NewPostBox onPost={addPost} />
      {posts.length === 0 ? (
        <div className="p-4 text-center text-gray-500">No posts yet</div>
      ) : (
        posts.map((post) => (
          <PostItem key={post.id} post={post} showMeta={true} emphasize={false} />
        ))
      )}
    </div>
  );
}
