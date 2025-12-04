import { ExamType } from "@/lib/api/types";
import {
	BookOpen,
	Calendar,
	FileText,
	NotebookPen,
	Trophy,
} from "lucide-react";

export const getExamTypeIcon = (type: ExamType) => {
	switch (type) {
		case "DAILY":
			return Calendar;
		case "MOCK":
			return NotebookPen;
		case "REVIEW":
			return BookOpen;
		case "FINAL":
			return Trophy;
		default:
			return FileText;
	}
};

export const getExamTypeColor = (type: ExamType) => {
	switch (type) {
		case "DAILY":
			return "text-green-600 bg-green-50";
		case "MOCK":
			return "text-purple-600 bg-purple-50";
		case "REVIEW":
			return "text-orange-600 bg-orange-50";
		case "FINAL":
			return "text-red-600 bg-red-50";
		default:
			return "text-gray-600 bg-gray-50";
	}
};

export const formatDate = (dateString: string) => {
	return new Date(dateString).toLocaleDateString("en-US", {
		month: "short",
		day: "numeric",
		year: "numeric",
		timeZone: "UTC",
	});
};

export const formatTime = (dateString: string) => {
	return (
		new Date(dateString).toLocaleTimeString("en-US", {
			hour: "2-digit",
			minute: "2-digit",
			timeZone: "UTC",
		}) + " UTC"
	);
};
