from langchain_openai import ChatOpenAI
from langchain.agents import AgentExecutor, create_react_agent
from langchain.agents import Tool

llm = ChatOpenAI(temperature=0,model_name='gpt-3.5-turbo')

import subprocess
from langchain.tools import Tool

class TerminalCallTool(Tool):
    def __init__(self):
        super().__init__(
            name="terminal_call",
            description="Executes terminal commands and returns the output.",
            func=self.run_terminal_command
        )

    def run_terminal_command(self, command: str) -> str:
        try:
            result = subprocess.run(command, shell=True, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            return result.stdout.decode('utf-8') + result.stderr.decode('utf-8')
        except subprocess.CalledProcessError as e:
            return f"An error occurred while executing the command: {e.stderr.decode('utf-8')}"

from langchain import hub

# https://smith.langchain.com/hub/hwchase17/react?organizationId=78d8038a-768c-5912-8e42-f162e94d62da
# zero shot
# prompt = hub.pull("hwchase17/react")

# https://smith.langchain.com/hub/hwchase17/react-chat?organizationId=78d8038a-768c-5912-8e42-f162e94d62da
# conversation 
prompt = hub.pull("hwchase17/react-chat")


agent = create_react_agent(
    llm=llm,
    tools=[TerminalCallTool()],
    prompt=prompt,
)

executor = AgentExecutor(
    tools=[TerminalCallTool()],
    llm=llm,
    agent=agent,
    verbose=True,
    handle_parsing_errors=True,
)

res = executor.invoke({"input": "Write a Golang hello world program in the current directory of my computer, and then execute it.","chat_history":""})
print(res)