import Link from 'next/link'
import { useRouter } from 'next/router'
import {useTheme} from 'next-themes'
import Modal from '../components/Modal'



export default function Layout({ children }) {
    const router = useRouter();
    const {theme, setTheme} = useTheme()
    const menuItems = [
        {
          href: '/dashboard/devices',
          title: 'Devices',
        },
        {
            href: '/dashboard/stats',
            title: 'Client Stats',
          },
        {
          href: '/dashboard/about',
          title: 'About',
        },
        
      ];
      
    return (
      
        
            <div className='min-h-screen flex flex-col justify-between overflow-hidden'>

              <div className="flex min-h-screen flex-col md:flex-row flex-1 overflow-hidden">
                <aside className='min-h-screen bg-red-100 w-full md:w-60 dark:bg-slate-600 overflow-hidden'>
                  <nav>
                    <ul>
                      {menuItems.map(({ href, title }) => (
                        <li className='m-2' key={title}>
                          <Link href={href} className={`flex p-2 bg-red-200 rounded hover:bg-red-400 cursor-pointer dark:bg-red-700 ${
                                router.asPath === href && 'bg-red-600 text-white dark:bg-slate-200'
                              }`}>

                              {title}
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </nav>
                </aside>
                <Modal/>
                <main className='flex-1 overflow-hidden'>{children}</main>
                
              </div>
            </div>
    );
  }