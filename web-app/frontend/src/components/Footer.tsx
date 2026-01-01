

function Footer() {
    return (
        <footer className="w-full py-4 bg-[#101010] border-t border-zinc-900 mt-auto">
            <div className="container mx-auto px-4 text-center text-[#bbbbbb] text-sm">
                &copy; 2025-{new Date().getFullYear()} Eichi Arakaki. All rights reserved.
            </div>
        </footer>
    );
}

export default Footer;