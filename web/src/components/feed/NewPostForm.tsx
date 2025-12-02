import type React from "react";
import { useRef, useState } from "react";

type Props = {
  onPost: (content: string) => void;
};

export default function NewPostForm() {
  const [text, setText] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!text.trim()) return;
    onPost(text);
    setText("");
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="border-b border-gray-300 p-3 flex items-start space-x-3 bg-white"
    >
      <div className="flex-1 overflow-hidden">
        <textarea
          ref={textareaRef}
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="What's happening?"
          className="w-full resize-none outline-none text-gray-800 placeholder-gray-400 bg-transparent min-h-[60px] overflow-hidden"
        />
        <div className="flex justify-end mt-2">
          <button
            type="submit"
            className="bg-indigo-600 text-white px-4 py-1.5 rounded-full text-sm font-semibold hover:bg-indigo-700 disabled:opacity-50"
            disabled={!text.trim()}
          >
            Post
          </button>
        </div>
      </div>
    </form>
  );
}
