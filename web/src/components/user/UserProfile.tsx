// If you're working on this please be familiar how react query works!!!
// https://tanstack.com/query/latest

import { useLoaderData } from "react-router";

export default function UserView() {
  const id: number = useLoaderData();
  return <></>;
}

// export default function UserView() {
//   const id: number = useLoaderData();
//   const { data, status, error } = useGetApiUsersId(id);
//
//   if (status === "pending") {
//     return <span>Loading...</span>;
//   }
//
//   if (status === "error") {
//     return <span>Error: {error.message}</span>;
//   }
//
//   if (status == "success") {
//     return (
//       <>
//         <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
//           <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
//             Profile
//           </header>
//           <div className="p-6">
//             <div className="flex flex-col items-center space-y-4">
//               <div className="avatar avatar-placeholder">
//                 <div className="bg-neutral text-neutral-content w-24 rounded-full">
//                   <span className="text-3xl">
//                     {data.data.username
//                       .replace(/^@/, "")
//                       .slice(0, 2)
//                       .toUpperCase()}
//                   </span>
//                 </div>
//               </div>
//               <div className="text-center">
//                 <div className="text-gray-600">{data.data.username}</div>
//                 <div className="text-2xl font-semibold text-gray-700">
//                   {data.data.bio}
//                 </div>
//                 <button
//                   type="button"
//                   onClick={() => { }}
//                   className="btn btn-primary"
//                 >
//                   Follow
//                 </button>
//               </div>
//             </div>
//           </div>
//         </div>
//         <div className="border-t border-gray-300" />
//       </>
//     );
//   }
// }
