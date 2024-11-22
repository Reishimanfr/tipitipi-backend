export default async function validateToken(setLoading : React.Dispatch<React.SetStateAction<boolean>>) {
  const token = localStorage.getItem("token");
  try {
    if (token === null) {
      console.debug("Token is invalid");
      return false;
    }
    const response = await fetch("http://localhost:2333/admin/validate", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.ok;
  } catch (error) {
    console.error("Token validation failed", error);
    return false;
  } finally {
    setLoading(false)
  }
}
