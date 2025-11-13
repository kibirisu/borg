export type User = {
  id: number;
  username: string;
  bio?: string;
  followers_count: number;
  following_count: number;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
  avatarColor?: string; // not in DB !!!
};
export type Post = {
  id: string;
  author: string;
  content: string;
  createdAt: string;
  likes?: number;
  replies?: number;
  reposts?: number;
};
export type Comment = {
  id: string;
  postId: string; // which post this comment belongs to
  author: string;
  content: string;
  createdAt: string;
};

export const samplePosts: Post[] = [
  {
    id: '1',
    author: '@ada',
    content: '# This\n is a comment, guys.',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
    likes: 12,
    replies: 3,
    reposts: 1,
  },
  {
    id: '2',
    author: '@grace',
    content: '### This is testing markdown\n`still testing markdown`\nThis is the realest comment that exists.',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
    likes: 45,
    replies: 8,
    reposts: 10,
  },
  {
    id: '3',
    author: '@jgrn',
    content: 'Hi\n I am a real person that is testing links\n[link to my page](/mypage).',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
    likes: 69,
    replies: 2,
    reposts: 1,
  },
  {
    id: '4',
    author: '@mzuck',
    content: 'Hi  \nI on the other hand am testing lists:  \n* a\n* b\n* c\n* last item',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
    likes: 45,
    replies: 8,
    reposts: 10,
  },
  {
    id: '5',
    author: '@ada',
    content: '# This\n is another post from ada.',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
    likes: 0,
    replies: 0,
    reposts: 0,
  },
];
export const sampleUsers: User[] = [
  {
    id: 1,
    username: '@ada',
    bio: 'Ada Lovelace — the first programmer.',
    followers_count: 120,
    following_count: 80,
    is_admin: false,
    created_at: '2025-10-12T09:15:32.000Z',
    updated_at: '2025-11-01T11:22:10.000Z',
    avatarColor: 'bg-purple-500',
  },
  {
    id: 2,
    username: '@grace',
    bio: 'Grace Hopper — pioneer of compilers.',
    followers_count: 200,
    following_count: 150,
    is_admin: false,
    created_at: '2025-10-11T08:42:18.000Z',
    updated_at: '2025-11-02T13:10:05.000Z',
    avatarColor: 'bg-green-500',
  },
  {
    id: 3,
    username: '@jgrn',
    bio: 'John Greene — likes writing and testing code.',
    followers_count: 75,
    following_count: 50,
    is_admin: false,
    created_at: '2025-10-15T10:55:44.000Z',
    updated_at: '2025-11-03T14:18:21.000Z',
    avatarColor: 'bg-red-500',
  },
  {
    id: 4,
    username: '@mzuck',
    bio: 'Mark Zuck — testing out lists and social stuff.',
    followers_count: 5000,
    following_count: 300,
    is_admin: true,
    created_at: '2025-09-20T12:00:00.000Z',
    updated_at: '2025-11-01T12:00:00.000Z',
    avatarColor: 'bg-blue-500',
  },
];

export const sampleComments: Comment[] = [
  // Comments for post 1 (Ada’s post)
  {
    id: 'c1',
    postId: '1',
    author: '@grace',
    content: 'Nice use of markdown, Ada! 💻',
    createdAt: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
  },
  {
    id: 'c2',
    postId: '1',
    author: '@jgrn',
    content: 'Love seeing some old-school markdown testing.',
    createdAt: new Date(Date.now() - 1000 * 60 * 25).toISOString(),
  },
  {
    id: 'c3',
    postId: '1',
    author: '@mzuck',
    content: 'Markdown forever! 😎',
    createdAt: new Date(Date.now() - 1000 * 60 * 20).toISOString(),
  },

  // Comments for post 2 (Grace’s post)
  {
    id: 'c4',
    postId: '2',
    author: '@ada',
    content: 'Looks great! Maybe add a code block too?',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 10).toISOString(),
  },
  {
    id: 'c5',
    postId: '2',
    author: '@jgrn',
    content: '`still testing markdown` — I see what you did there 😄',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 9).toISOString(),
  },
  {
    id: 'c6',
    postId: '2',
    author: '@mzuck',
    content: 'This is indeed the *realest* comment.',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 8).toISOString(),
  },

  // Comments for post 3 (John’s post)
  {
    id: 'c7',
    postId: '3',
    author: '@ada',
    content: 'Checked your link — works perfectly!',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
  },
  {
    id: 'c8',
    postId: '3',
    author: '@grace',
    content: 'Links are essential for a connected web! 🌐',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 1.5).toISOString(),
  },

  // Comments for post 4 (Mark’s post)
  {
    id: 'c9',
    postId: '4',
    author: '@jgrn',
    content: 'Lists are always satisfying to read.',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 5).toISOString(),
  },
  {
    id: 'c10',
    postId: '4',
    author: '@ada',
    content: 'You forgot the bullet for item D 😅',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 4.5).toISOString(),
  },
  {
    id: 'c11',
    postId: '4',
    author: '@grace',
    content: 'Classic markdown list — looks perfect to me!',
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 4).toISOString(),
  },

  // Comments for post 5 (Ada’s second post)
  {
    id: 'c12',
    postId: '5',
    author: '@grace',
    content: 'Another Ada original — keep them coming!',
    createdAt: new Date(Date.now() - 1000 * 60 * 10).toISOString(),
  },
];
