import { MoveRight } from "lucide-react";
import { Link, useNavigate } from "react-router";
import { useState } from "react";
import axios from "axios";
import { toast, ToastContainer } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

const Register = () => {
    const navigate = useNavigate();
    

    const [email, setEmail] = useState("");
    const [fullName, setFullName] = useState("");
    const [phone, setPhone] = useState("");
    const [password, setPassword] = useState("");


    const submitHandle = async (e) => {
        e.preventDefault();

        try {
            const response = await axios.post(`${import.meta.env.VITE_BACK_END}/v1/create_user`, {
                user_name: email,
                full_name: fullName,
                password,
                phone_number: phone,
            });

            if (response.status === 200 && response.data.user_id) {
                toast.success("Registration successful!");
                setTimeout(() => navigate("/auth/login"), 1500);
            } else {
                toast.error("Registration failed. Please try again.");
            }
        } catch (error) {
            console.error("Register error:", error);
            toast.error(error?.response?.data?.message || "An error occurred during registration.");
        }
    };

    return (
        <div className="lg:container mx-auto p-[80px]">
            <div className="max-w-[648px] w-full min-h-[382px] p-[31px] mx-auto flex items-center justify-center flex-col rounded-lg border-[1px] border-[#9a9caa]">
                <h3 className="text-3xl text-[#272343] font-semibold font-inter mb-5 capitalize">Register</h3>

                <form onSubmit={submitHandle} className="flex flex-col items-center w-full space-y-4">
                    <input
                        type="text"
                        placeholder="Your Name..."
                        className="w-full h-[50px] bg-[#f0f2f3] rounded-lg pl-3.5"
                        value={fullName}
                        onChange={(e) => setFullName(e.target.value)}
                    />
                    <input
                        type="email"
                        placeholder="Your Email..."
                        className="w-full h-[50px] bg-[#f0f2f3] rounded-lg pl-3.5"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                    <input
                        type="text"
                        placeholder="Your Phone..."
                        className="w-full h-[50px] bg-[#f0f2f3] rounded-lg pl-3.5"
                        value={phone}
                        onChange={(e) => setPhone(e.target.value)}
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
                        className="w-full h-[50px] bg-[#bbf5ed] rounded-lg pl-3.5 text-base text-white font-semibold font-inter capitalize flex items-center justify-center cursor-pointer gap-2.5"
                    >
                        Register <MoveRight />
                    </button>
                </form>
                
                <p className="text-base text-[#272343] font-normal font-inter flex items-center justify-center gap-2.5 mt-4">
                    Don't have an account? 
                    <Link to={'/auth/login'} className="text-[#007580]">Login</Link>
                </p>
            </div>

            {/* Toast notifications container */}
            <ToastContainer />
        </div>
    );
};

export default Register;
