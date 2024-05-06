export type ResultStats = {
	time: number;
	total: number;
	index: number;
	pages: number;
	count: number;
	params: object;
	facet: Map<string, FacetCount[]> | null;
};

export type FacetCount = {
	label: string;
	count: number;
};

export type ResultResponse<T> = {
	stats: ResultStats;
	items: T[];
	message: string | null;
};

export type SearchProblemResult = ResultResponse<Problem>;

export type Problem = {
	problemId: string;
	problemTitle: string;
	problemUrl: string;
	contestId: string;
	contestTitle: string;
	contestUrl: string;
	difficulty: number | null;
	color: string | null;
	startAt: string;
	duration: number;
	rateChange: string;
	category: string;
};

export type SearchUserResult = ResultResponse<User>;

export type User = {
	userId: string;
	rating: number;
	highestRating: number;
	affiliation: string | null;
	birthYear: number | null;
	country: string | null;
	crown: string | null;
	joinCount: number;
	rank: number;
	activeRank: number | null;
	wins: number;
	color: string;
	userUrl: string;
};

export type SearchSubmissionResult = ResultResponse<Submission>;

export type Submission = {
	submissionId: number;
	submittedAt: string;
	submissionUrl: string;
	problemId: string;
	problemTitle: string;
	contestId: string;
	contestTitle: string;
	category: string;
	difficulty: number;
	color: string;
	userId: string;
	language: string;
	point: number;
	length: number;
	result: string;
	executionTime: number | null;
};
