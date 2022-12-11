import React from 'react';
import Logo from "../components/Logo";
import {useTheme} from "next-themes";
import{SunIcon ,MoonIcon} from "@heroicons/react/24/solid";
import { useState,useEffect } from 'react';

    const Header = () => {
        const [mounted, setMounted] = useState(false);

      const {systemTheme , theme, setTheme} = useTheme ();
      useEffect(() =>{
        setMounted(true);
      },[])

      const renderThemeChanger= () => {
        if(!mounted) return null;

          const currentTheme = theme === "system" ? systemTheme : theme ;

          if(currentTheme ==="dark"){
            return (
              <SunIcon className=" w-10 h-10 text-yellow-500 " role="button" onClick={() => setTheme('light')} />
            )
          }

          else {
            return (
              <MoonIcon className=" w-10 h-10 text-gray-900 " role="button" onClick={() => setTheme('dark')} />
            )
          }
       };

      return (
        <header className="border-slate-700 fixed inset-x-0 top-0 border bg-red-200 h-15 shadow-sm dark:bg-slate-700 border-black">
          <div className=" container px-4 sm:px-6 py-4 flex justify-between items-center">
            {/* Logo */}
            <Logo />
            <div className='hidden w-full lg:inline-flex lg:flex-grow lg:w-auto'>
          <div className='lg:inline-flex lg:flex-row lg:ml-auto lg:w-auto w-full lg:items-center items-start'>
               
                {renderThemeChanger()}
          </div>
            
           </div> 
          </div>
        </header>
      );
    };

export default Header;
