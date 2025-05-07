import { Navigate, useLocation } from "react-router";

const AuthCheck = ({ children }) => {
    const token = localStorage.getItem("token");
    const isAuthenticated = !!token; 
    
    const location = useLocation();

    if (!isAuthenticated) {
        return <Navigate to="/auth/login" state={{ from: location }} />;
    }

    return children;
};

export default AuthCheck;
