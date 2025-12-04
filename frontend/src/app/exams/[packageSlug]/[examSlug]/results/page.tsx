"use client";

import { withAuth } from "@/contexts/AuthContext";
import { useParams, useSearchParams } from "next/navigation";
import { ExamResultInterface } from "./components";

function ExamResultPage() {
	const params = useParams();
	const searchParams = useSearchParams();
	const packageSlug = params.packageSlug as string;
	const examSlug = params.examSlug as string;
	const sessionId = searchParams.get("session");

	return (
		<ExamResultInterface
			packageSlug={packageSlug}
			examSlug={examSlug}
			sessionId={sessionId}
		/>
	);
}

export default withAuth(ExamResultPage);
