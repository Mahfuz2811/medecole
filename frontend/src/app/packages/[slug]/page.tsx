import { PackageDetailClient } from "./PackageDetailClient";

interface PackageDetailPageProps {
	params: Promise<{
		slug: string;
	}>;
}

export default async function PackageDetailPage({
	params,
}: PackageDetailPageProps) {
	const { slug } = await params;
	return <PackageDetailClient slug={slug} />;
}
