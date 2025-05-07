import { MoveRight } from "lucide-react";
import { Link, useNavigate } from "react-router";
import { useState } from "react";
import axios from "axios";
import { toast, ToastContainer } from "react-toastify"; 

import 'react-toastify/dist/ReactToastify.css';

const Login = () => {
    const navigate = useNavigate();
    
    const [user_name, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const submitHandle = async (e) => {
        e.preventDefault();

        try {
            const response = await axios.post(`${process.env.VITE_BACK_END}/v1/login`, {
                user_name,
                password,
            });
            if (response.status === 200 && response.data.access_token ) {
                localStorage.setItem("id",response.data.user_id)
                localStorage.setItem("token", response.data.access_token );

                toast.success("Login successful!");
                navigate("/");
            } else {
                toast.error("Login failed! Please check your credentials."); 
            }
        } catch (error) {
            console.error("Error during API request:", error);
            toast.error("Error during login! Please try again later."); 
        }
    };

    return (
        <div className="lg:container mx-auto p-[80px]">
            <div className="max-w-[648px] w-full min-h-[382px] p-[31px] mx-auto flex items-center justify-center flex-col rounded-lg border-[1px] border-[#9a9caa]">
                <h3 className="text-3xl text-[#272343] font-semibold font-inter mb-5 capitalize">Login</h3>

                <form onSubmit={submitHandle} className="flex flex-col items-center w-full space-y-4">
                    <input
                        type="email"
                        placeholder="Your Email..."
                        className="w-full h-[50px] bg-[#f0f2f3] rounded-lg pl-3.5"
                        value={user_name}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                    <input
                        type="password"
                        placeholder="Your Password..."
                        className="w-full h-[50px] bg-[#f0f2f3] rounded-lg pl-3.5"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <button
                        type="submit"
                        className="w-full h-[50px] bg-[#007580] rounded-lg pl-3.5 text-base text-white font-semibold font-inter capitalize flex items-center justify-center cursor-pointer gap-2.5"
                    >
                        Login <MoveRight />
                    </button>
                </form>

                <p className="text-base text-[#272343] font-normal font-inter flex items-center justify-center gap-2.5 mt-4">
                    Don't have an account? 
                    <Link to={'/auth/register'} className="text-[#007580]">Register</Link>
                </p>
            </div>
            
            {/* Toast notifications container */}
            <ToastContainer />
        </div>
    );
};

export default Login;
