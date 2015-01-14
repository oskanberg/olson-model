#!/usr/bin/python

import pygame
from pygame.locals import *
import time, sys, json

white = (255,255,255)
black = (0,0,0)
red = (255, 0, 0)


class Pane(object):
    def __init__(self):
        pygame.init()
        pygame.display.set_caption('Simulation')
        self.screen = pygame.display.set_mode((1000,1000), 0, 32)
        self.screen.fill((white))
        pygame.display.flip()


    def drawAgent(self, agent):
    	colour = black
    	renderFitness = False
    	if agent['AgentType'] == 'Predator':
    		colour = red
    		renderFitness = True

    	x = agent['Position']['Location']['X']
    	y = agent['Position']['Location']['Y']
        self.rect = pygame.draw.rect(self.screen, (colour), (x , y, 10, 10), 2)

        if renderFitness:
        	myfont = pygame.font.SysFont("monospace", 15)
	        label = myfont.render(str(agent['Fitness']), 1, black)
	        self.screen.blit(label, (x, y))

    def load_and_draw(self):
    	json_data = {}
    	with open(sys.argv[1], 'r') as f:
    		json_data = json.load(f)
    	for step in json_data['Steps']:
	    	for position in step['Positions']:
	    		self.drawAgent(position)
	    	pygame.display.flip()
	    	time.sleep(0.001)
	    	self.screen.fill(white)	
    	raw_input()



display = Pane()
display.load_and_draw()