// import { useEffect , useState} from "react";
// import Unauthorized from "../pages/errorPages/unauthorized";

// const ProtectedRoute = (props: any) => {
//   const [isAuthorized,setIsAuthorized] = useState(false)
//   useEffect(() => {
//     console.log("checking authorization");
//     const token = localStorage.getItem("token");

//     if (token === null) {
//       console.debug("Token is invalid");
//       //setIsAuthorized(false);
//       return;
//     }

//     validateToken(token);
//     console.log("isAuthorized: " + isAuthorized);
//   }) , [isAuthorized];

//   async function validateToken(token: string) {
//     try {
//       const response = await fetch("http://localhost:2333/admin/validate", {
//         method: "POST",
//         headers: { Authorization: `Bearer ${token}` },
//       });
//       setIsAuthorized(response.ok);
//     } catch (error) {
//       console.error("Token validation failed", error);
//       setIsAuthorized(false);
//     }
//   }

//   if (!isAuthorized) {
//     // Jeśli użytkownik nie jest autoryzowany, przekieruj na stronę logowania
//     return <Unauthorized/>
//   }
//   // Jeśli użytkownik jest autoryzowany, pokaż chronioną stronę
//   return props.element;
// };

// export default ProtectedRoute;
