import { useState } from 'react';
import { useLoaderData } from 'react-router';
import { useQuery } from '@tanstack/react-query';
import type { components } from '../../lib/api/v1';
import { type Post } from './feedData';
import NewPostBox from './NewPostBox';
import PostItem from './PostItem';
import type { Client } from '../../lib/api/client';

export const loader = (client: Client) => async () => {
  const opts = client.$api.queryOptions('get', '/api/posts', {});
  try {
    await client.queryClient.ensureQueryData(opts);
  } catch (error) {
    // If backend is not ready, return opts anyway - component will handle retry
    console.warn('Failed to prefetch posts:', error);
  }
  return { opts };
};

export default function MainFeed() {
  const { opts } = useLoaderData() as Awaited<ReturnType<ReturnType<typeof loader>>>;
  const { data: apiPosts, isPending, error } = useQuery(opts);
  const [localPosts, setLocalPosts] = useState<Post[]>([]);

  // Map API posts to feed format
  const posts: Post[] = apiPosts
    ? apiPosts.map((p: components['schemas']['Post']) => {
        // Handle createdAt - it might be a string or Date object
        let createdAtStr: string;
        if (p.createdAt instanceof Date) {
          createdAtStr = p.createdAt.toISOString();
        } else if (typeof p.createdAt === 'string') {
          createdAtStr = p.createdAt;
        } else {
          createdAtStr = new Date().toISOString();
        }
        
        return {
          id: String(p.id),
          author: p.username ? `@${p.username}` : `@user${p.userID}`,
          content: p.content,
          createdAt: createdAtStr,
          likes: p.likeCount,
          replies: p.commentCount,
          reposts: p.shareCount,
          userID: p.userID,
        };
      })
    : [];

  // Combine with local posts (newly created)
  const allPosts = [...localPosts, ...posts];

  function addPost(content: string) {
    const newPost: Post = {
      id: String(Date.now()),
      author: '@you',
      content: content.trim(),
      createdAt: new Date().toISOString(),
      likes: 0,
      replies: 0,
      reposts: 0,
    };
    setLocalPosts((p) => [newPost, ...p]);
  }

  if (isPending) {
    return (
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Home
        </header>
        <div className="p-4">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Home
        </header>
        <div className="p-4 text-red-600">
          Error loading posts. Please make sure the backend server is running on port 8080.
          <br />
          <small>{error instanceof Error ? error.message : String(error)}</small>
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
      {allPosts.map((post) => (
        <PostItem key={post.id} post={post} />
      ))}
    </div>
  );
}
