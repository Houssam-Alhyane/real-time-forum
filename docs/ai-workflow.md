# Dual-Tool AI Workflow Plan

## Phase 1: The Specification (Human)

- Your action: Write a 2–5 sentence description of your feature or fix in plain English.
- Include explicit constraints such as endpoints, database choices, limits, or other requirements.

## Phase 2: The Planning Step (Gemini Browser)

- Your action: Open Gemini in your browser and submit your specification using Prompt 1.

- Gemini's output: An ordered subtask list detailing the steps and exact file paths to modify.
- The plan must not include code.

## Phase 3: The Plan Review (Human)

- Your action: Read Gemini's plan in the browser.
- Correct any bad architectural assumptions or missed edge cases in the chat.
- Keep refining until the plan is perfect.
- Copy the final approved plan.

## Phase 4: The Execution Step (VS Code Copilot / Codex)

- Your action: Open your codebase in VS Code.
- Provide your AI tool with the full context of the approved plan.
- Command it to execute only Subtask 1 using Prompt 2.
- VS Code AI output: The exact code for that specific subtask.

## Phase 5: Code Review & Integration (Human + VS Code AI)

- Your action: Review the code line-by-line.
- If needed, use Copilot to explain or review the generated snippet.
- Test it and commit the change.
- Move to the next subtask using Prompt 3.

## The Modified Prompts

### Prompt 1: For Gemini Browser (Planning Phase)

> "I want you to act strictly as an architectural planning assistant. Break down the following specification into an ordered list of subtasks, specifying the exact file paths that need to be created or modified. Do not write any code, function skeletons, or implementation details. Provide only the plain English step-by-step plan.
>
> Here is the specification:
> [INSERT YOUR 2-5 SENTENCE SPECIFICATION HERE]"

### Prompt 2: For VS Code Copilot / Codex (First Execution Step)

> "I am using a strict 'separate planning from execution' workflow. Here is the overall approved architectural plan for this feature:
>
> [PASTE THE ENTIRE PLAN FROM GEMINI HERE]
>
> We are executing this one step at a time. Write the code exclusively for the current subtask: '[INSERT SUBTASK 1 HERE]'. Do not write code for any other subtasks yet."
> Constraints: only pick what is for the exact mentionned subtask from raw-requirements.md exclude the other features constraints

### Prompt 3: For VS Code Copilot / Codex (Subsequent Execution Steps)

> "The previous step is integrated. Now, write the code exclusively for the next subtask from our plan: '[INSERT THE NEXT SUBTASK HERE]'. Do not move ahead to the rest of the list."
> Constraints: only pick what is for the exact mentionned subtask from raw-requirements.md exclude the other features constraints

## Prompt Master (Optional)

Use this prompt when you need to break a large coding project into Epics before planning.

> "I need to break down a massive coding project into distinct, manageable Epics so I can plan and execute them one at a time. Read the master specification below.
> important: you must completely skip already implemented features
> Your task is to act as an Agile Product Owner. Group the requirements into logical Epics. For each Epic, provide:
>
> - A short title.
> - A 2-3 sentence 'Plain English Specification' defining the end state of that specific Epic (ignoring implementation details).
> - The dependencies (e.g., 'Requires Epic 1 to be finished').
>
> Do not write any code or file-level subtasks yet. Just give me the high-level Epic breakdown.
>
> [PASTE YOUR FULL PROJECT SPECIFICATION HERE]"

Once Gemini gives you that list, pick the first Epic, take its 2-3 sentence specification, and feed it back into Gemini to generate the subtask plan with file paths, as outlined in the blog post.
