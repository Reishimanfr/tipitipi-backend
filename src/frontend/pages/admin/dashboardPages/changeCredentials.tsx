import { useState ,useEffect } from "react";
import validateToken from "../../../components/validate";
import Unauthorized from "../../errorPages/unauthorized";

const ChangeCredentials = () => {
  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);
  useEffect(() => {
    const ValidateAuthorization = async () => {
      setIsAuthorized(await validateToken(setLoading));
    };
    ValidateAuthorization();
  }, []);
  if (loading) {
    return <div>Loading</div>;
  }
  if (!isAuthorized) {
    return <Unauthorized />;
  }
  return <div></div>;
};

export default ChangeCredentials;
