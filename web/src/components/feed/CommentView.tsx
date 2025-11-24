import { UserX } from 'lucide-react';
import { useParams } from 'react-router';
import { type Comment, type Post, sampleComments, samplePosts } from '../feed/feedData';
import PostItem from '../feed/PostItem';
import TopAppBar from '../TopAppBar';

/**
 * View a single post (enlarged) and display its comments below.
 */
export default function CommentView({ onTopBarSearch }: any) {
  const params = useParams();
  const postId = params.postId;

  // Pretend to fetch the post and comments
  const post = samplePosts.find((p: Post) => p.id === postId);
  const comments = sampleComments.filter((c: Comment) => c.postId === postId);

  if (!post) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen text-center p-6">
        <div className="flex flex-col items-center gap-4 max-w-md">
          <div className="bg-red-100 text-red-600 p-4 rounded-full">
            <UserX className="w-10 h-10" />
          </div>
          <h1 className="text-3xl font-bold">Sorry, that post does not exist</h1>
          <p className="text-gray-500">The post you’re looking for might have been deleted or never existed.</p>
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

  return (
    <div>
      <TopAppBar onSearch={onTopBarSearch} />
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Post
        </header>

        {/* The main post (enlarged) */}
        <div className="border-b border-gray-200 p-4 bg-gray-50">
          <PostItem post={post} showMeta={true} emphasize={true} />
        </div>

        {/* Comments list */}
        <section className="divide-y divide-gray-200">
          {comments.length === 0 ? (
            <div className="p-6 text-gray-500 text-center">No comments yet. Be the first to reply!</div>
          ) : (
            comments.map((comment) => (
              <div key={comment.id} className="p-4">
                <PostItem post={comment} showMeta={false} />
              </div>
            ))
          )}
        </section>
      </div>
    </div>
  );
}
