const Footer = () => {
    return (
      <footer className="mt-auto bg-red-200 py-6 border-slate-700 border fixed inset-x-0 bottom-0 p-4 dark:bg-slate-700">
        <div className="bg-red-200 text-center text-sm text-gray-500 dark:bg-slate-700">
          <span className="dark:text-gray-100 text-red-900 font-bold text-lg mr-2">WaaS</span>   &copy; {new Date().getFullYear()}   All Rights Reserved
        </div>
      </footer>
    );
  };

export default Footer;
