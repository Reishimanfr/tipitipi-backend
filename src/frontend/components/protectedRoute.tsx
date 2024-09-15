import { useNavigate } from 'react-router-dom';

const ProtectedRoute = (props : any) => {
    const navigate = useNavigate()
  if (!props.isAuthorized) {
    // Jeśli użytkownik nie jest autoryzowany, przekieruj na stronę logowania
    navigate("/admin/login")
  }
  // Jeśli użytkownik jest autoryzowany, pokaż chronioną stronę
  return props.children;
};

export default ProtectedRoute;
