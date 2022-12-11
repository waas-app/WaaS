import { LightningBoltIcon } from "@heroicons/react/24/solid";
import Link from "next/link";
import Image from "next/image";
import logo from "../image/wireguard.png"

const Logo = () => {
    return (
      
        <Link href="/" className="my-2 flex items-center space-x-1 text-red-900">
         
          <span className="dark:text-gray-100 text-black-900 font-bold text-lg font-bold text-3xl font-sans tracking-tight whitespace-nowrap">WaaS</span>
          <Image class="float-right ..." src={logo} width={50} height ={50}/>
          
        </Link>
    
    )
}

export default Logo;
