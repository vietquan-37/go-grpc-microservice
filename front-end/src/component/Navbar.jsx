import { House } from "lucide-react";
import { Link, useNavigate } from "react-router"; 

function Navbar() {
    const navigate = useNavigate(); 

    const logout = () => {
        localStorage.clear()
        navigate("/auth/login"); 
    };
    const token = localStorage.getItem("token");
    return (
        <div className="flex items-center bg-[#272343] h-[45px] w-full">
            <div className="lg:container flex items-center w-full">
                <div className="hidden sm:flex items-center">
                    <p className="flex text-sm font-normal text-white capitalize ml-1 items-center gap-0.5">
                        <House /> Simple fe
                    </p>
                </div>
                <div className="ml-auto flex gap-6">
                    {!token ? (
                        <>
                            <button className="text-sm text-white font-inter font-normal capitalize">
                                <Link to="/auth/login">Login</Link>
                            </button>
                            <button className="text-sm text-white font-inter font-normal capitalize">
                                <Link to="/auth/register">Register</Link>
                            </button>
                        </>
                    ) : (
                    
                        <button
                            className="text-sm text-white font-inter font-normal capitalize"
                            onClick={logout}
                        >
                            Logout
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}

export default Navbar;
