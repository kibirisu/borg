import { Heart, MessageCircle, Repeat, Share2 } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { Link } from 'react-router';
import type { Post } from './feedData';

type Props = {
  post: Post;
  showMeta: boolean;
  emphasize: boolean;
};

export function renderMetaButtons(showMeta: boolean, post: Post) {
  if (showMeta) {
    return (
      <div className="flex justify-between mt-3 text-gray-500 text-sm max-w-md">
        <button type="button" className="flex items-center space-x-1 hover:text-blue-500 transition">
          <MessageCircle size={16} /> <span>{post.replies}</span>
        </button>
        <button type="button" className="flex items-center space-x-1 hover:text-green-500 transition">
          <Repeat size={16} /> <span>{post.reposts}</span>
        </button>
        <button type="button" className="flex items-center space-x-1 hover:text-pink-500 transition">
          <Heart size={16} /> <span>{post.likes}</span>
        </button>
        <button type="button" className="flex items-center space-x-1 hover:text-gray-700 transition">
          <Share2 size={16} />
        </button>
      </div>
    );
  } else {
    return;
  }
}
export default function PostItem({ post, showMeta = true, emphasize = false }: Props) {
  const metaButtons = renderMetaButtons(showMeta, post);
  return (
    <div className="border-b border-gray-200 p-4 hover:bg-gray-50 transition-colors">
      <div className="flex space-x-3">
        <div className="flex-1">
          <div className="flex items-center space-x-1 mb-2">
            <Link
              to={`/profile/${post.author.replace(/^@/, '')}`}
              className={`hover:underline ${emphasize ? 'font-bold text-xl text-gray-900' : 'font-semibold text-gray-900'}`}
            >
              {post.author}
            </Link>
          </div>
          <div className={`prose max-w-none ${emphasize ? 'text-lg text-gray-900' : 'text-gray-800'}`}>
            <ReactMarkdown>{post.content}</ReactMarkdown>
          </div>
          {metaButtons}
        </div>
      </div>
    </div>
  );
}
