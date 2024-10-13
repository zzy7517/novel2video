import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export const showToast = (message: string) => {
    toast(message, {
        position: "top-center",
        autoClose: 100,
        closeOnClick: false,
        draggable: false,
        closeButton: false,
    });
};