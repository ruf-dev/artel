import {Path} from "@/app/routing/Router.tsx"
import cls from "@/segments/Topbar/components/TopbarBrand/TopbarBrand.module.css"
import BrandMarkIcon from "@/segments/Topbar/components/BrandMarkIcon/BrandMarkIcon.tsx"
export default function TopbarBrand() {
    return (
        <a className={cls.Brand} href={Path.HomePage} aria-label="Artel home">
            <BrandMarkIcon />
            <span className={cls.BrandWord}>artel</span>
        </a>
    )
}

