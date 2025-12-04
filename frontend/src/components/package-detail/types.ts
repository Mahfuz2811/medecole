import { PackageResponse } from "@/lib/api/types";

export interface PackageDetailProps {
	package: PackageResponse;
}

export interface PackageActionProps extends PackageDetailProps {
	isAuthenticated: boolean;
	onShowLogin: () => void;
}

export type MobileTab = "overview" | "exams" | "enroll";

export interface MobileTabNavigationProps extends PackageDetailProps {
	activeTab: MobileTab;
	onTabChange: (tab: MobileTab) => void;
}

export interface LoginModalProps extends PackageDetailProps {
	onClose: () => void;
}
