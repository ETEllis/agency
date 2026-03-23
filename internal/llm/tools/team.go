package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ETEllis/teamcode/internal/team"
)

type TeamTool struct {
	service *team.Service
}

func NewTeamTool() *TeamTool {
	return &TeamTool{
		service: team.NewService(),
	}
}

func (t *TeamTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "team_create_context",
		Description: "Create or update team context with charter, roles, and goals. Call this first to establish a team.",
		Parameters: map[string]any{
			"team_name": "string - Name of the team",
			"charter":   "string - Team's purpose and mission",
			"roles":     "map - Role definitions (architect, implementer, reviewer, integrator)",
		},
		Required: []string{"team_name"},
	}
}

func (t *TeamTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, ok := params["team_name"].(string)
	if !ok || teamName == "" {
		return NewTextErrorResponse("team_name is required"), nil
	}

	charter, _ := params["charter"].(string)

	tc, err := t.service.Context.CreateContext(ctx, teamName, charter, nil)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to create team context: %s", err)), nil
	}

	return WithResponseMetadata(NewTextResponse(fmt.Sprintf("Team context created for %s", teamName)), map[string]any{
		"team_name": tc.TeamName,
		"charter":   tc.Charter,
	}), nil
}

type TeamAddRoleTool struct {
	service *team.Service
}

func NewTeamAddRoleTool() *TeamAddRoleTool {
	return &TeamAddRoleTool{
		service: team.NewService(),
	}
}

func (t *TeamAddRoleTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "team_add_role",
		Description: "Add a role definition to the team context",
		Parameters: map[string]any{
			"team_name":     "string - Name of the team",
			"role_name":     "string - Name of the role (e.g., architect, implementer, reviewer)",
			"responsible":   "string - What this role is responsible for",
			"current_focus": "string - Current focus of this role (optional)",
		},
		Required: []string{"team_name", "role_name", "responsible"},
	}
}

func (t *TeamAddRoleTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, _ := params["team_name"].(string)
	roleName, _ := params["role_name"].(string)
	responsible, _ := params["responsible"].(string)
	currentFocus, _ := params["current_focus"].(string)

	role := team.Role{
		Name:        roleName,
		Responsible: responsible,
	}
	if currentFocus != "" {
		role.CurrentFocus = currentFocus
	}

	_, err := t.service.Context.AddRole(ctx, teamName, roleName, role)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to add role: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("Role '%s' added to team %s", roleName, teamName)), nil
}

type TeamAssignRoleTool struct {
	service *team.Service
}

func NewTeamAssignRoleTool() *TeamAssignRoleTool {
	return &TeamAssignRoleTool{
		service: team.NewService(),
	}
}

func (t *TeamAssignRoleTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "team_assign_role",
		Description: "Assign an agent to a role in the team",
		Parameters: map[string]any{
			"team_name":  "string - Name of the team",
			"role_name":  "string - Name of the role",
			"agent_name": "string - Name of the agent to assign",
		},
		Required: []string{"team_name", "role_name", "agent_name"},
	}
}

func (t *TeamAssignRoleTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, _ := params["team_name"].(string)
	roleName, _ := params["role_name"].(string)
	agentName, _ := params["agent_name"].(string)

	_, err := t.service.Context.AssignAgentToRole(ctx, teamName, roleName, agentName)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to assign agent to role: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("Agent '%s' assigned to role '%s' in team %s", agentName, roleName, teamName)), nil
}

type TaskBoardTool struct {
	service *team.Service
}

func NewTaskBoardTool() *TaskBoardTool {
	return &TaskBoardTool{
		service: team.NewService(),
	}
}

func (t *TaskBoardTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "task_create",
		Description: "Create a task on the team's task board",
		Parameters: map[string]any{
			"team_name": "string - Name of the team",
			"task_id":   "string - Unique identifier for the task",
			"column":    "string - Column to add task to (backlog, ready, in_progress, in_review, done, blocked)",
		},
		Required: []string{"team_name", "task_id"},
	}
}

func (t *TaskBoardTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, _ := params["team_name"].(string)
	taskID, _ := params["task_id"].(string)
	column, _ := params["column"].(string)
	if column == "" {
		column = "backlog"
	}

	_, err := t.service.Board.AddTaskToColumn(ctx, teamName, taskID, column)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to add task: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("Task '%s' added to %s column in team %s", taskID, column, teamName)), nil
}

type TaskMoveTool struct {
	service *team.Service
}

func NewTaskMoveTool() *TaskMoveTool {
	return &TaskMoveTool{
		service: team.NewService(),
	}
}

func (t *TaskMoveTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "task_move",
		Description: "Move a task between columns on the task board",
		Parameters: map[string]any{
			"team_name":   "string - Name of the team",
			"task_id":     "string - Task identifier",
			"from_column": "string - Source column",
			"to_column":   "string - Destination column",
			"agent":       "string - Agent to assign (optional, for in_progress)",
		},
		Required: []string{"team_name", "task_id", "from_column", "to_column"},
	}
}

func (t *TaskMoveTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, _ := params["team_name"].(string)
	taskID, _ := params["task_id"].(string)
	fromColumn, _ := params["from_column"].(string)
	toColumn, _ := params["to_column"].(string)
	agent, _ := params["agent"].(string)

	_, err := t.service.Board.MoveTask(ctx, teamName, taskID, fromColumn, toColumn, agent)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to move task: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("Task '%s' moved from %s to %s", taskID, fromColumn, toColumn)), nil
}

type HandoffCreateTool struct {
	service *team.Service
}

func NewHandoffCreateTool() *HandoffCreateTool {
	return &HandoffCreateTool{
		service: team.NewService(),
	}
}

func (t *HandoffCreateTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "handoff_create",
		Description: "Create a handoff to transfer a completed task to another agent",
		Parameters: map[string]any{
			"team_name":    "string - Name of the team",
			"task_id":      "string - Task being handed off",
			"from_agent":   "string - Agent handing off",
			"to_agent":     "string - Agent receiving the handoff",
			"work_summary": "string - Summary of completed work (optional)",
			"artifacts":    "array - Files created by this work (optional)",
		},
		Required: []string{"team_name", "task_id", "from_agent", "to_agent"},
	}
}

func (t *HandoffCreateTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, _ := params["team_name"].(string)
	taskID, _ := params["task_id"].(string)
	fromAgent, _ := params["from_agent"].(string)
	toAgent, _ := params["to_agent"].(string)
	workSummary, _ := params["work_summary"].(string)

	h, err := t.service.Handoff.Create(ctx, teamName, taskID, fromAgent, toAgent, workSummary, nil, nil)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to create handoff: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("Handoff created: %s (task: %s, from: %s, to: %s)", h.ID, h.TaskID, h.FromAgent, h.ToAgent)), nil
}

type HandoffAcceptTool struct {
	service *team.Service
}

func NewHandoffAcceptTool() *HandoffAcceptTool {
	return &HandoffAcceptTool{
		service: team.NewService(),
	}
}

func (t *HandoffAcceptTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "handoff_accept",
		Description: "Accept a pending handoff",
		Parameters: map[string]any{
			"team_name":   "string - Name of the team",
			"handoff_id":  "string - Handoff identifier",
			"accepted_by": "string - Agent accepting the handoff",
		},
		Required: []string{"team_name", "handoff_id", "accepted_by"},
	}
}

func (t *HandoffAcceptTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, _ := params["team_name"].(string)
	handoffID, _ := params["handoff_id"].(string)
	acceptedBy, _ := params["accepted_by"].(string)

	h, err := t.service.Handoff.Accept(ctx, teamName, handoffID, acceptedBy)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to accept handoff: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("Handoff %s accepted by %s", h.ID, acceptedBy)), nil
}

type InboxReadTool struct {
	service *team.Service
}

func NewInboxReadTool() *InboxReadTool {
	return &InboxReadTool{
		service: team.NewService(),
	}
}

func (t *InboxReadTool) Info() ToolInfo {
	return ToolInfo{
		Name:        "inbox_read",
		Description: "Read messages from your inbox",
		Parameters: map[string]any{
			"team_name":   "string - Name of the team",
			"agent_name":  "string - Your agent name",
			"unread_only": "bool - Only show unread messages (optional, default false)",
		},
		Required: []string{"team_name", "agent_name"},
	}
}

func (t *InboxReadTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	teamName, _ := params["team_name"].(string)
	agentName, _ := params["agent_name"].(string)
	unreadOnly, _ := params["unread_only"].(bool)

	msgs, err := t.service.Inbox.ReadInbox(ctx, teamName, agentName, unreadOnly)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to read inbox: %s", err)), nil
	}

	if len(msgs) == 0 {
		return NewTextResponse("No messages"), nil
	}

	var response string
	for _, msg := range msgs {
		response += fmt.Sprintf("[%s] From %s: %s\n", msg.Timestamp, msg.From_, msg.Text)
	}

	return NewTextResponse(response), nil
}
