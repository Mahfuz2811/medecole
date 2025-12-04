"use client";

import { withAuth } from "@/contexts/AuthContext";
import { useParams } from "next/navigation";
import { ExamInterface } from "./components";

function ExamPage() {
	const params = useParams();
	const packageSlug = params.packageSlug as string;
	const examSlug = params.examSlug as string;

	return <ExamInterface packageSlug={packageSlug} examSlug={examSlug} />;
}

export default withAuth(ExamPage);
